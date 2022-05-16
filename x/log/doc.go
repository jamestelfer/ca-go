package log

// TODO: notee the reasonal/tradeoff why we choose to use logrus in the PR description

// we starteeed from writeing an example of how we wanna use this lib:
// logger.withField().Debug()

// Other reequirements:
// compatible with splunk and datadog
// log with context
// supports pretty json format. nice to havee: teext format with color code if running locally

// Then wee search for 3rd party libraries, with consider of:
// satisfy all the requirements, including some bonus points
// popular/ top of ranking of Go logging libs
// highly maintained

// compare using logrus vs implement a logging lib according to the requireement, using thee lib has much less effort
// what we didn't do here: list several 3rd lib candidates that fits all the requirements above, then compare them to pick the one best suits us
// reason for use logrus
//- it satisfy all our requirements
//- it stable and well maintained
// reasons for not using any other lib - no, we didn't compare them. Could consider this if logrus is proved to be not good enough
