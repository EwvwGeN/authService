package template

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_HappyPass(t *testing.T) {
	err := PrepareTemplates()
	require.NoError(t, err)

	link := "testedlink"
	expectedTmpl := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body>
    <div>
        <p>Confirm the registration by clicking on the button</p>
        <a href="%s"><button>Confirm</button></a>
        </br>
        <p>If the button does not work, follow the link:</p>
        <p>%s</p>
    </div>
</body>
</html>`, link, link)

	body, err := Register(link)
	require.NoError(t, err)
	require.Equal(t, expectedTmpl, string(body))
}