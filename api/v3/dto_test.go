package v3

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorSchema(t *testing.T) {
	t.Parallel()

	t.Run("bot error", func(t *testing.T) {
		t.Parallel()

		jsonData, err := os.ReadFile("testdata/error_bot.json")
		require.NoError(t, err)

		var errorSchema ErrorSchema
		err = json.Unmarshal(jsonData, &errorSchema)
		require.NoError(t, err)

		expected := ErrorSchema{
			Data: ErrorDataSchema{
				Error:     "Looks like you are a bot. Please, pass captcha or obtain proper valid API key",
				ErrorCode: 9000,
			},
			Result: "error",
		}
		assert.Equal(t, expected, errorSchema)

		actualJSON, err := json.Marshal(&errorSchema)
		require.NoError(t, err)

		assert.JSONEq(t, string(jsonData), string(actualJSON))
	})

	t.Run("missing params error", func(t *testing.T) {
		t.Parallel()

		jsonData, err := os.ReadFile("testdata/error_missing_params.json")
		require.NoError(t, err)

		var errorSchema ErrorSchema
		err = json.Unmarshal(jsonData, &errorSchema)
		require.NoError(t, err)

		expected := ErrorSchema{
			Data: ErrorDataSchema{
				Error:     "Missing parameters: ['query']",
				ErrorCode: 103,
			},
			Result: "error",
		}
		assert.Equal(t, expected, errorSchema)

		actualJSON, err := json.Marshal(&errorSchema)
		require.NoError(t, err)

		assert.JSONEq(t, string(jsonData), string(actualJSON))
	})
}

func TestSearchResponseSchema(t *testing.T) {
	t.Parallel()

	t.Run("basic", func(t *testing.T) {
		t.Parallel()

		// this test is failing
		// the reason is known and the fix is on the way
		t.Skip("the test fails due to API inconsistency, skip for now")

		jsonData, err := os.ReadFile("testdata/search_response1.json")
		require.NoError(t, err)

		var response SearchResponseSchema
		err = json.Unmarshal(jsonData, &response)
		require.NoError(t, err)

		expected := SearchResponseSchema{
			Data:   response.Data,
			Result: "OK",
		}
		assert.Equal(t, expected, response)

		actualJSON, err := json.Marshal(&response)
		require.NoError(t, err)

		assert.JSONEq(t, string(jsonData), string(actualJSON))

		data, err := response.Data.AsSearchResponseDataSchema()
		require.NoError(t, err)

		actualDataJSON, err := json.Marshal(&data)
		require.NoError(t, err)

		assert.JSONEq(t, string(response.Data.union), string(actualDataJSON))
	})

	t.Run("search results and fields", func(t *testing.T) {
		t.Parallel()

		jsonData, err := os.ReadFile("testdata/search_response1.json")
		require.NoError(t, err)

		var response SearchResponseSchema
		err = json.Unmarshal(jsonData, &response)
		require.NoError(t, err)

		data, err := response.Data.AsSearchResponseDataSchema()
		require.NoError(t, err)

		assert.Len(t, data.Search, 20)
		assert.Equal(t, data.Total, 339)
		assert.Equal(t, data.MaxSearchSize, 100)
	})
}

func TestSearchByIDResponseSchema(t *testing.T) {
	t.Parallel()

	t.Run("with references", func(t *testing.T) {
		t.Parallel()

		jsonData, err := os.ReadFile("testdata/search_by_id_response1.json")
		require.NoError(t, err)

		var response SearchByIDResponseSchema
		err = json.Unmarshal(jsonData, &response)
		require.NoError(t, err)

		expected := SearchByIDResponseSchema{
			Data:   response.Data,
			Result: "OK",
		}
		assert.Equal(t, expected, response)

		actualJSON, err := json.Marshal(&response)
		require.NoError(t, err)

		assert.JSONEq(t, string(jsonData), string(actualJSON))

		data, err := response.Data.AsSearchByIDResponseDataSchema()
		require.NoError(t, err)

		actualDataJSON, err := json.Marshal(&data)
		require.NoError(t, err)

		assert.JSONEq(t, string(response.Data.union), string(actualDataJSON))
	})

	t.Run("without references", func(t *testing.T) {
		t.Parallel()

		jsonData, err := os.ReadFile("testdata/search_by_id_response2.json")
		require.NoError(t, err)

		var response SearchByIDResponseSchema
		err = json.Unmarshal(jsonData, &response)
		require.NoError(t, err)

		expected := SearchByIDResponseSchema{
			Data:   response.Data,
			Result: "OK",
		}
		assert.Equal(t, expected, response)

		actualJSON, err := json.Marshal(&response)
		require.NoError(t, err)

		assert.JSONEq(t, string(jsonData), string(actualJSON))

		data, err := response.Data.AsSearchByIDResponseDataSchema()
		require.NoError(t, err)

		actualDataJSON, err := json.Marshal(&data)
		require.NoError(t, err)

		assert.JSONEq(t, string(response.Data.union), string(actualDataJSON))
	})
}
