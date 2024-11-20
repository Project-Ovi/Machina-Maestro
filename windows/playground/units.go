package playground

type measurementSystem int

const metricSystem = measurementSystem(1)
const imperialSystem = measurementSystem(2)
const SISystem = measurementSystem(3)

/*
Units is used to define the measuring system used
1: Metric
2: Imperial
3: SI
*/
var thisMeasurementSystem = metricSystem

func (sys measurementSystem) Name() string {
	switch sys {
	case metricSystem:
		return "Metric"
	case imperialSystem:
		return "Imperial"
	case SISystem:
		return "SI"
	}
	return "Unknown system"
}
