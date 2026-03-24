package logging

// Test helpers — these symbols are only visible to the logging_test package.

var NewAsyncWriterForTest = newAsyncWriter
var NewSlogHandlerForTest = newSlogHandler
var NewSinkWriterForTest = newSinkWriter

// Done returns the done channel of the AsyncWriter, allowing tests to wait for shutdown.
func (w *AsyncWriter) Done() <-chan struct{} { return w.done }
