package fixtures

import (
	"github.com/DataDog/raclette/model"
)

var defaultAggregators = []string{"service", "name", "resource"}

// TestStatsBucket returns a fixed stats bucket to be used in unit tests
func TestStatsBucket() model.StatsBucket {
	sb := model.NewStatsBucket(0, 1e9)
	sb.HandleSpan(TestSpan(), defaultAggregators)
	return sb
}

// StatsBucketWithSpans returns a stats bucket populated with spans stats
func StatsBucketWithSpans(s []model.Span) model.StatsBucket {
	sb := model.NewStatsBucket(0, 1e9)
	for _, s := range s {
		sb.HandleSpan(s, defaultAggregators)
	}
	return sb
}

// RandomStatsBucket returns a bucket made from n random spans, useful to run benchmarks and tests
func RandomStatsBucket(n int) model.StatsBucket {
	spans := make([]model.Span, 0, n)
	for i := 0; i < n; i++ {
		spans = append(spans, RandomSpan())
	}

	return StatsBucketWithSpans(spans)
}

// TestDistroValues is a pre-defined list of values
var TestDistroValues = []int64{
	49873, 81744, 46545, 43680, 7535, 33334, 93009, 23777, 33471, 68629,
	94601, 83827, 3556, 15913, 84957, 368, 71879, 73687, 55039, 89704,
	98733, 40820, 62839, 26673, 55731, 45477, 15893, 45488, 72297, 29134,
	57683, 6782, 10496, 16713, 62976, 7545, 87884, 7963, 16105, 28633,
	19613, 33881, 53049, 39639, 68647, 99105, 95954, 79172, 65798, 32334,
	66448, 13783, 56688, 17350, 42414, 18336, 63655, 59545, 42014, 74478,
	70263, 6860, 19339, 36375, 72034, 51899, 98473, 22231, 57126, 4482,
	31985, 35335, 89732, 58843, 28695, 50653, 23740, 29245, 72152, 16566,
	19598, 94928, 88210, 9813, 61112, 14225, 282, 40069, 80421, 71429,
	30896, 38353, 34031, 65116, 20348, 57019, 91726, 3143, 48396, 25658,
	465, 48299, 4127, 73883, 99755, 95259, 79187, 59794, 25740, 62633,
	61585, 26320, 96966, 57059, 2201, 20065, 58359, 75706, 67622, 90459,
	19300, 40384, 98456, 65224, 15020, 35819, 48079, 8554, 41658, 22967,
	28764, 78538, 78314, 73160, 8707, 83916, 7982, 38096, 45418, 78655,
	27987, 41748, 84730, 91216, 83098, 49090, 426, 48221, 26862, 70959,
	32132, 19862, 95997, 3027, 19438, 38393, 33338, 25567, 14618, 31610,
	88956, 4252, 81845, 77757, 58023, 64701, 24762, 11909, 79436, 67507,
	63004, 62749, 55296, 88204, 43255, 15385, 4404, 75079, 32425, 32088,
	35378, 83907, 15201, 37043, 49320, 3941, 10696, 77039, 45697, 33241,
	6414, 91211, 11473, 39560, 1833, 59542, 30878, 40429, 18136, 45348,
	20395, 18976, 22945, 72978, 11297, 49834, 74443, 32954, 87079, 43619,
	51680, 44241, 24348, 30395, 8241, 6038, 34042, 10788, 43017, 1706,
	41296, 87732, 17445, 90738, 12690, 7810, 14243, 10162, 26128, 36418,
	90821, 63677, 87168, 35589, 89271, 2882, 19680, 46951, 67143, 99086,
	40945, 88011, 88062, 8742, 24121, 41593, 52634, 94285, 84646, 46255,
	66570, 11781, 4395, 82956, 98527, 6198, 10414, 71817, 52338, 8849,
	70229, 54649, 98215, 81781, 28883, 50424, 65524, 89666, 18922, 25075,
	26313, 91007, 45330, 52683, 19222, 58549, 15102, 66637, 11874, 96489,
	20224, 96151, 38772, 77736, 26639, 63909, 5960, 81147, 68183, 15503,
	98095, 45086, 79831, 95974, 69140, 38202, 40126, 96299, 48670, 29259,
	21494, 60618, 45045, 63612, 56271, 57411, 90412, 43692, 4981, 79404,
	11842, 39727, 12257, 90435, 6909, 61222, 34525, 45393, 39051, 45634,
	11202, 86878, 89570, 84142, 8400, 30596, 55909, 22552, 45053, 34014,
	3546, 41567, 54300, 233, 53248, 78597, 39224, 87627, 82985, 43282,
	37318, 11994, 47289, 6375, 14274, 74678, 7444, 64063, 95054, 94864,
	56093, 11942, 66802, 71928, 816, 13229, 62403, 78549, 41223, 55717,
	19609, 56257, 28648, 6162, 3943, 9800, 97273, 30486, 50528, 66419,
	56069, 77098, 99676, 50095, 25915, 5126, 88303, 91216, 39747, 35313,
	67128, 33430, 80861, 4598, 98636, 58579, 464, 58865, 54999, 2770,
	50827, 31275, 28270, 81736, 50019, 30829, 7715, 28098, 59506, 93275,
	59696, 3620, 78626, 94467, 99199, 56480, 81559, 66099, 14158, 14121,
	58014, 77264, 36713, 23639, 99892, 28986, 15902, 71818, 40326, 6597,
	66142, 63904, 12735, 23989, 43671, 45438, 76740, 41381, 61377, 240,
	15913, 96435, 68748, 14924, 73254, 86370, 37633, 61430, 99398, 45688,
	8955, 8474, 97979, 39943, 93195, 65534, 22004, 19573, 53598, 14585,
	36601, 99530, 91841, 44689, 63644, 84307, 72608, 78387, 8859, 78854,
	50002, 22510, 85289, 95122, 5656, 25727, 79150, 55133, 4004, 96902,
	29830, 77912, 85867, 90171, 82337, 44654, 96195, 59459, 5902, 91724,
	67780, 7250, 85047, 34558, 38288, 78736, 19084, 714, 67720, 72898,
	48739, 61426, 557, 25487, 59289, 47253, 63439, 37965, 81039, 5683,
	7797, 32679, 23594, 65206, 19993, 25043, 25180, 20326, 23150, 36051,
	64304, 40757, 57203, 26517, 68184, 67824, 95437, 75023, 88923, 70288,
	24445, 3, 95502, 77711, 56441, 7932, 67526, 68888, 99420, 55438,
	46474, 56435, 72679, 99497, 49292, 97114, 74148, 11560, 90975, 11458,
	41169, 3235, 69486, 98718, 79108, 91634, 55222, 55298, 8990, 86267,
	64122, 69275, 50964, 52229, 95153, 90588, 80232, 32330, 76329, 21423,
	67743, 58663, 78473, 63279, 990, 37566, 14986, 86231, 85598, 48049,
	10363, 57368, 31711, 8906, 21830, 80262, 95792, 17164, 60127, 57617,
	20080, 21982, 64448, 20778, 72023, 86362, 36221, 55531, 23085, 99240,
	67901, 90321, 20114, 62605, 96437, 24478, 53523, 28354, 80996, 80790,
	88883, 15785, 91293, 77907, 90565, 68434, 38138, 38726, 89991, 71803,
	63103, 77849, 170, 30055, 2028, 16229, 41089, 11047, 43713, 45225,
	13700, 47201, 6036, 12316, 99542, 53145, 79478, 36265, 3113, 10984,
	49406, 60035, 62615, 80977, 71344, 14200, 95778, 22538, 60343, 67009,
	63429, 32294, 27237, 68984, 12944, 32231, 55999, 37897, 90091, 80466,
	95801, 65865, 96564, 66561, 31327, 52672, 71584, 2776, 579, 91374,
	55089, 78267, 77595, 83646, 243, 58118, 79231, 99188, 62236, 44332,
	81093, 38651, 73028, 99672, 68818, 9953, 93758, 93236, 97302, 92746,
	33019, 14922, 29229, 54180, 52829, 90520, 38644, 51461, 29513, 66800,
	22806, 54867, 48009, 46546, 25875, 30956, 31243, 68299, 16312, 85165,
	64305, 77372, 10692, 6157, 59324, 29112, 63886, 72133, 85611, 38971,
	38992, 44689, 24522, 24774, 73909, 24398, 6723, 43141, 30123, 70649,
	56382, 67159, 26385, 65003, 98672, 69931, 66304, 1286, 23984, 7956,
	37911, 53510, 43011, 75474, 73917, 1584, 66755, 64636, 14254, 74482,
	21556, 59100, 17851, 55708, 22718, 24043, 74123, 40832, 2753, 86226,
	24531, 75018, 42006, 96396, 32645, 66235, 68342, 21044, 33145, 56726,
	96180, 31556, 54782, 93490, 74270, 71721, 33153, 73493, 71298, 35333,
	22479, 57204, 38705, 94009, 44158, 15787, 26768, 64158, 22096, 88571,
	45340, 91679, 24695, 64220, 25843, 70219, 66173, 81020, 89491, 74995,
	41378, 69249, 58966, 57816, 99462, 13518, 79224, 43384, 48332, 54966,
	86665, 31802, 46214, 42160, 11404, 63242, 36740, 20740, 56772, 34968,
	2681, 93064, 73128, 93007, 29572, 29621, 55075, 8266, 91077, 94482,
	95789, 60063, 85533, 50362, 41231, 49712, 35175, 75798, 9857, 66660,
	9971, 63999, 8772, 30798, 76458, 82119, 11532, 67945, 18010, 18487,
	58024, 27592, 45792, 25909, 37864, 35835, 38423, 77954, 78700, 89430,
	50421, 59400, 76022, 95436, 65119, 27730, 66893, 1156, 85816, 88373,
	50148, 25239, 59895, 60979, 87409, 99146, 3145, 44179, 51101, 26657,
	62651, 32522, 16767, 50925, 70089, 57291, 1766, 20794, 20314, 13281,
	72751, 37028, 28871, 92784, 42352, 58356, 41326, 85658, 8758, 44966,
	11090, 2746, 71165, 58595, 52442, 38558, 74826, 6675, 32401, 83481,
	93010, 19397, 86634, 3022, 94732, 11983, 84975, 44749, 80986, 5166,
	78078, 52211, 82787, 56617, 61960, 23051, 64815, 49412, 16945, 10430,
	61046, 38824, 85281, 59365, 87997, 96782, 4978, 29164, 68042, 66505,
	68244, 92104, 2331, 10527, 15497, 57600, 91716, 12689, 51048, 95748,
	79084, 26550, 93068, 89239, 32297, 30275, 3483, 19284, 83045, 45812,
	47572, 89241, 23722, 22646, 54408, 35989, 39531, 25405, 81469, 69026,
	59956, 88882, 47029, 32217, 9265, 9337, 15567, 71576, 82557, 83448,
	76538, 95379, 97595, 30781, 54709, 40266, 97288, 89581, 97335, 54606,
	64572, 99834, 97581, 10704, 51460, 54803, 41618, 41760, 31663, 42939,
	10327, 63265, 48904, 79260, 26562, 11528, 97745, 78918, 94479, 19453,
}

// TestDistribution returns a distribution with pre-defined values
func TestDistribution() model.Distribution {
	tgs := model.NewTagsFromString("service:X,name:Y,host:Z")
	d := model.NewDistribution("duration", "duration|service:X,name:Y,host:Z", tgs)
	for i, v := range TestDistroValues {
		d.Add(float64(v), uint64(i))
	}

	return d
}
