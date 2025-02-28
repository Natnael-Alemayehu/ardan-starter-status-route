package mid

import "github.com/Natnael-Alemayehu/chat_clone/chat/foundation/web"

// isError tests if the Encoder has an error inside of it.
func isError(e web.Encoder) error {
	err, isError := e.(error)
	if isError {
		return err
	}
	return nil
}
