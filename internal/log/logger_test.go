package log

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/Pallinder/go-randomdata"
)

func newFields(initial Fields) Fields {
	fields := Fields{}
	for k, v := range initial {
		fields[k] = v
	}
	for i := 0; i < randomdata.Number(1, 5); i++ {
		name := "field" + strconv.Itoa(i)
		fields[name] = randomdata.SillyName()
	}
	return fields
}

var requiredServiceFields = Fields{
	FieldCategory: CategoryApplication,
	FieldService:  randomdata.SillyName() + "-service",
}

type logWriterFuncs struct {
	Write, WriteLine func(...interface{})
	WriteFormat      func(string, ...interface{})
}

func getLogWriterFuncs(level Level) (*logWriterFuncs, error) {
	switch level {
	case LevelDebug:
		return &logWriterFuncs{Write: Debug, WriteFormat: Debugf, WriteLine: Debugln}, nil
	case LevelInfo:
		return &logWriterFuncs{Write: Info, WriteFormat: Infof, WriteLine: Infoln}, nil
	case LevelWarn:
		return &logWriterFuncs{Write: Warn, WriteFormat: Warnf, WriteLine: Warnln}, nil
	case LevelError:
		return &logWriterFuncs{Write: Error, WriteFormat: Errorf, WriteLine: Errorln}, nil
	case LevelFatal:
		return &logWriterFuncs{Write: Fatal, WriteFormat: Fatalf, WriteLine: Fatalln}, nil
	case LevelPanic:
		return &logWriterFuncs{Write: Panic, WriteFormat: Panicf, WriteLine: Panicln}, nil
	}
	return nil, fmt.Errorf("unhandled level: %s", level)
}

var _ = Describe("Logging Benchmarks", func() {
	var (
		err    error
		logger Logger
	)

	BeforeEach(func() {
		fields := Fields{
			FieldCategory: CategoryApplication,
			FieldService:  "test-service",
		}
		logger, err = NewServiceLogger(fields, LevelDebug)
		if err != nil {
			Fail(fmt.Sprintf("setup benchmark: %s", err))
		}
	})

	Measure("it should log efficiently", func(b Benchmarker) {
		var out bytes.Buffer
		var message = "This is a test message containing about 50 characters."
		// There is one output destination per runtime.  Do not run parallel benchmarks.
		SetOutput(&out)
		runtime := b.Time("runtime", func() {
			logger.Debug(message)
		})

		Ω(runtime.Seconds()).Should(BeNumerically("<", 0.005), "Logging a message took too long")
	}, 10)
})

var _ = Describe("Logging Tests", func() {
	var (
		err    error
		logger Logger
	)

	Context("ServiceLogger", func() {
		DescribeTable("required fields",
			func(fields Fields, expectError bool) {
				logger, err = NewServiceLogger(fields, LevelDebug)
				if expectError {
					Expect(err).ShouldNot(Succeed())
				} else {
					Expect(err).Should(Succeed())
				}
			},
			Entry("missing FieldCategory", Fields{FieldService: "x"}, true),
			Entry("missing FieldService", Fields{FieldCategory: "x"}, true),
			Entry("none missing", requiredServiceFields, false),
		)

		DescribeTable("filters log messages",
			func(minLevel Level, logLevel Level) {
				var out bytes.Buffer
				SetOutput(&out)
				if logger, err = NewServiceLogger(requiredServiceFields, minLevel); err != nil {
					Fail(fmt.Sprintf("create logger: %s", err))
				}
				SetLogger(logger)
				var funcs *logWriterFuncs
				if funcs, err = getLogWriterFuncs(logLevel); err != nil {
					Fail(fmt.Sprintf("get log-write funcs: %s", err))
				}
				funcs.Write()
				Expect(out.Len()).Should(Equal(0))
				funcs.WriteFormat("")
				Expect(out.Len()).Should(Equal(0))
				funcs.WriteLine()
				Expect(out.Len()).Should(Equal(0))
			},
			Entry("info filter > debug entry", LevelInfo, LevelDebug),
			Entry("warn filter > info entry", LevelWarn, LevelInfo),
			Entry("error filter > warn entry", LevelError, LevelWarn),
			Entry("fatal filter > error entry", LevelFatal, LevelError),
			Entry("panic filter > fatal entry", LevelPanic, LevelError),
		)

		DescribeTable("does not filter debug/info/warn/error log messages",
			func(minLevel Level, logLevel Level) {
				var out bytes.Buffer
				SetOutput(&out)
				if logger, err = NewServiceLogger(requiredServiceFields, minLevel); err != nil {
					Fail(fmt.Sprintf("create logger: %s", err))
				}
				SetLogger(logger)
				var funcs *logWriterFuncs
				if funcs, err = getLogWriterFuncs(logLevel); err != nil {
					Fail(fmt.Sprintf("get log-write funcs: %s", err))
				}
				funcs.Write()
				Expect(out.Len()).Should(BeNumerically(">", 0))
				funcs.WriteFormat("")
				Expect(out.Len()).Should(BeNumerically(">", 0))
				funcs.WriteLine()
				Expect(out.Len()).Should(BeNumerically(">", 0))
			},
			Entry("debug filter <= debug entry", LevelDebug, LevelDebug),
			Entry("info filter <= info entry", LevelInfo, LevelInfo),
			Entry("warn filter <= warn entry", LevelWarn, LevelWarn),
			Entry("error filter <= error entry", LevelError, LevelError),
		)

		Describe("logs fields to output destination", func() {
			var (
				out          bytes.Buffer
				loggedFields Fields
			)
			BeforeEach(func() {
				SetOutput(&out)
				if logger, err = NewServiceLogger(requiredServiceFields, LevelDebug); err != nil {
					Fail(fmt.Sprintf("create logger: %s", err))
				}
				SetLogger(logger)
				loggedFields = newFields(requiredServiceFields)
			})
			JustBeforeEach(func() {
				WithFields(loggedFields).Debug()
			})
			It("should succeed", func() {
				for k, v := range loggedFields {
					Expect(out.String()).Should(ContainSubstring(`"%s":"%s"`, k, v))
				}
			})
		})

		/* There is no way to intercept os.Exit().

		Describe("exits on fatal level", func() {
		})
		*/

		Describe("causes a panic on LevelPanic", func() {
			JustBeforeEach(func() {
				if logger, err = NewServiceLogger(requiredServiceFields, LevelDebug); err != nil {
					Fail(fmt.Sprintf("create logger: %s", err))
				}
			})
			It("should panic", func() {
				loggerPanic := func() { logger.Panic() }
				Ω(loggerPanic).Should(GomegaPanic())
			})
		})
	})

	Describe("DefaultLogger", func() {
		var (
			out bytes.Buffer
		)

		Describe("is available after package init", func() {
			BeforeEach(func() {
				SetOutput(&out)
			})

			JustBeforeEach(func() {
				Warn("hello")
			})

			It("should log output", func() {
				Ω(out.String()).Should(ContainSubstring("hello"))
			})
		})

		Describe("is configurable", func() {
			BeforeEach(func() {
				SetOutput(&out)
				fields := Fields{
					FieldCategory: CategoryApplication,
					FieldService:  "test-service",
				}
				logger, err := NewServiceLogger(fields, LevelDebug)
				if err != nil {
					Fail(fmt.Sprintf("create service logger: %s", err))
				}
				SetLogger(logger)
			})

			JustBeforeEach(func() {
				Debugf("%s", "goodbye")
			})

			It("should log output", func() {
				Ω(out.String()).Should(ContainSubstring("goodbye"))
			})
		})
	})
})
