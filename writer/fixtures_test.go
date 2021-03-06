package writer

import (
	"math/rand"
	"sync"

	log "github.com/cihub/seelog"

	"github.com/DataDog/datadog-trace-agent/fixtures"
)

// payloadConstructedHandlerArgs encodes the arguments passed to a PayloadConstructedHandler call.
type payloadConstructedHandlerArgs struct {
	payload *Payload
	stats   interface{}
}

// testEndpoint represents a mocked endpoint that replies with a configurable error and records successful and failed
// payloads.
type testEndpoint struct {
	sync.RWMutex
	err             error
	successPayloads []Payload
	errorPayloads   []Payload
}

// Write mocks the writing of a payload to a remote endpoint, recording it and replying with the configured error (or
// success in its absence).
func (e *testEndpoint) Write(payload *Payload) error {
	e.Lock()
	defer e.Unlock()
	if e.err != nil {
		e.errorPayloads = append(e.errorPayloads, *payload)
	} else {
		e.successPayloads = append(e.successPayloads, *payload)
	}
	return e.err
}

func (e *testEndpoint) Error() error {
	e.RLock()
	defer e.RUnlock()
	return e.err
}

// ErrorPayloads returns all the error payloads registered with the test endpoint.
func (e *testEndpoint) ErrorPayloads() []Payload {
	e.RLock()
	defer e.RUnlock()
	return e.errorPayloads
}

// SuccessPayloads returns all the success payloads registered with the test endpoint.
func (e *testEndpoint) SuccessPayloads() []Payload {
	e.RLock()
	defer e.RUnlock()
	return e.successPayloads
}

// SetError sets the passed error on the endpoint.
func (e *testEndpoint) SetError(err error) {
	e.Lock()
	defer e.Unlock()
	e.err = err
}

func (e *testEndpoint) String() string {
	return "testEndpoint"
}

// RandomPayload creates a new payload instance using random data and up to 32 bytes.
func RandomPayload() *Payload {
	return RandomSizedPayload(rand.Intn(32))
}

// RandomSizedPayload creates a new payload instance using random data with the specified size.
func RandomSizedPayload(size int) *Payload {
	return NewPayload(fixtures.RandomSizedBytes(size), fixtures.RandomStringMap())
}

// testPayloadSender is a PayloadSender that is connected to a testEndpoint, used for testing.
type testPayloadSender struct {
	testEndpoint *testEndpoint
	BasePayloadSender
}

// newTestPayloadSender creates a new instance of a testPayloadSender.
func newTestPayloadSender() *testPayloadSender {
	testEndpoint := &testEndpoint{}
	return &testPayloadSender{
		testEndpoint:      testEndpoint,
		BasePayloadSender: *NewBasePayloadSender(testEndpoint),
	}
}

// Start asynchronously starts this payload sender.
func (c *testPayloadSender) Start() {
	go c.Run()
}

// Run executes the core loop of this sender.
func (c *testPayloadSender) Run() {
	c.exitWG.Add(1)
	defer c.exitWG.Done()

	for {
		select {
		case payload := <-c.in:
			stats, err := c.send(payload)

			if err != nil {
				c.notifyError(payload, err, stats)
			} else {
				c.notifySuccess(payload, stats)
			}
		case <-c.exit:
			return
		}
	}
}

// Payloads allows access to all payloads recorded as being successfully sent by this sender.
func (c *testPayloadSender) Payloads() []Payload {
	return c.testEndpoint.SuccessPayloads()
}

// Endpoint allows access to the underlying testEndpoint.
func (c *testPayloadSender) Endpoint() *testEndpoint {
	return c.testEndpoint
}

func (c *testPayloadSender) setEndpoint(endpoint Endpoint) {
	c.testEndpoint = endpoint.(*testEndpoint)
}

// testPayloadSenderMonitor monitors a PayloadSender and stores all events
type testPayloadSenderMonitor struct {
	SuccessEvents []SenderSuccessEvent
	FailureEvents []SenderFailureEvent
	RetryEvents   []SenderRetryEvent

	sender PayloadSender

	exit   chan struct{}
	exitWG sync.WaitGroup
}

// newTestPayloadSenderMonitor creates a new testPayloadSenderMonitor monitoring the specified sender.
func newTestPayloadSenderMonitor(sender PayloadSender) *testPayloadSenderMonitor {
	return &testPayloadSenderMonitor{
		sender: sender,
		exit:   make(chan struct{}),
	}
}

// Start asynchronously starts this payload monitor.
func (m *testPayloadSenderMonitor) Start() {
	go m.Run()
}

// Run executes the core loop of this monitor.
func (m *testPayloadSenderMonitor) Run() {
	m.exitWG.Add(1)
	defer m.exitWG.Done()

	for {
		select {
		case event := <-m.sender.Monitor():
			if event == nil {
				continue
			}

			switch event := event.(type) {
			case SenderSuccessEvent:
				m.SuccessEvents = append(m.SuccessEvents, event)
			case SenderFailureEvent:
				m.FailureEvents = append(m.FailureEvents, event)
			case SenderRetryEvent:
				m.RetryEvents = append(m.RetryEvents, event)
			default:
				log.Errorf("Unknown event of type %T", event)
			}
		case <-m.exit:
			return
		}
	}
}

// Stop stops this payload monitor and waits for it to stop.
func (m *testPayloadSenderMonitor) Stop() {
	close(m.exit)
	m.exitWG.Wait()
}

// SuccessPayloads returns a slice containing all successful payloads.
func (m *testPayloadSenderMonitor) SuccessPayloads() []Payload {
	result := make([]Payload, len(m.SuccessEvents))

	for i, successEvent := range m.SuccessEvents {
		result[i] = *successEvent.Payload
	}

	return result
}

// FailurePayloads returns a slice containing all failed payloads.
func (m *testPayloadSenderMonitor) FailurePayloads() []Payload {
	result := make([]Payload, len(m.FailureEvents))

	for i, successEvent := range m.FailureEvents {
		result[i] = *successEvent.Payload
	}

	return result
}

// RetryPayloads returns a slice containing all failed payloads.
func (m *testPayloadSenderMonitor) RetryPayloads() []Payload {
	result := make([]Payload, len(m.RetryEvents))

	for i, successEvent := range m.RetryEvents {
		result[i] = *successEvent.Payload
	}

	return result
}
