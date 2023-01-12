package valuegeneration

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

// getTestABIArguments obtains ABI arguments of various types for use in testing ABI value related methods.
func getTestABIArguments() abi.Arguments {
	// Define our argument types to test round trip serialization for.
	args := abi.Arguments{
		{
			Name: "testAddress",
			Type: abi.Type{
				Elem:          nil,
				Size:          20,
				T:             abi.AddressTy,
				TupleRawName:  "",
				TupleElems:    nil,
				TupleRawNames: nil,
				TupleType:     nil,
			},
			Indexed: false,
		},
		{
			Name: "testString",
			Type: abi.Type{
				Elem:          nil,
				Size:          0,
				T:             abi.StringTy,
				TupleRawName:  "",
				TupleElems:    nil,
				TupleRawNames: nil,
				TupleType:     nil,
			},
			Indexed: false,
		},
		{
			Name: "testDynamicBytes",
			Type: abi.Type{
				Elem:          nil,
				Size:          0,
				T:             abi.BytesTy,
				TupleRawName:  "",
				TupleElems:    nil,
				TupleRawNames: nil,
				TupleType:     nil,
			},
			Indexed: false,
		},
		{
			Name: "testBool",
			Type: abi.Type{
				Elem:          nil,
				Size:          0,
				T:             abi.BoolTy,
				TupleRawName:  "",
				TupleElems:    nil,
				TupleRawNames: nil,
				TupleType:     nil,
			},
			Indexed: false,
		},
	}

	// Append all fixed byte sizes
	for i := 1; i <= 32; i++ {
		args = append(args, abi.Argument{
			Name: fmt.Sprintf("testBytes%d", i),
			Type: abi.Type{
				Elem:          nil,
				Size:          i,
				T:             abi.FixedBytesTy,
				TupleRawName:  "",
				TupleElems:    nil,
				TupleRawNames: nil,
				TupleType:     nil,
			},
			Indexed: false,
		})
	}

	// Append all integer types
	for i := 8; i <= 256; i += 8 {
		// Add our signed/unsigned types
		args = append(args,
			abi.Argument{
				Name: fmt.Sprintf("int%d", i),
				Type: abi.Type{
					Elem:          nil,
					Size:          i,
					T:             abi.IntTy,
					TupleRawName:  "",
					TupleElems:    nil,
					TupleRawNames: nil,
					TupleType:     nil,
				},
				Indexed: false,
			},
			abi.Argument{
				Name: fmt.Sprintf("uint%d", i),
				Type: abi.Type{
					Elem:          nil,
					Size:          i,
					T:             abi.UintTy,
					TupleRawName:  "",
					TupleElems:    nil,
					TupleRawNames: nil,
					TupleType:     nil,
				},
				Indexed: false,
			})
	}

	// Save our slice of the arguments with basic types
	basicArgs := args[:]

	// Add arguments that are arrays of each basic type
	for _, basicArg := range basicArgs {
		// Alias our arg so when we get a pointer to it, we don't cause a memory leak in this loop
		basicArg := basicArg

		// Add a slice/array of this basic type.
		args = append(args,
			abi.Argument{
				Name: fmt.Sprintf("testSlice (%v)", basicArg.Type.GetType().String()),
				Type: abi.Type{
					Elem:          &basicArg.Type,
					Size:          0,
					T:             abi.SliceTy,
					TupleRawName:  "",
					TupleElems:    nil,
					TupleRawNames: nil,
					TupleType:     nil,
				},
				Indexed: false,
			},
			abi.Argument{
				Name: fmt.Sprintf("testArray (%v)", basicArg.Type.GetType().String()),
				Type: abi.Type{
					Elem:          &basicArg.Type,
					Size:          5,
					T:             abi.ArrayTy,
					TupleRawName:  "",
					TupleElems:    nil,
					TupleRawNames: nil,
					TupleType:     nil,
				},
				Indexed: false,
			},
		)
	}

	// TODO: Add tuple argument.
	return args
}

// TestABIRoundtripEncodingAllTypes runs tests to ensure ABI value encoding works round-trip for argument values of all
// types. It generates values using a ValueGenerator, then encodes them, decodes them, and re-encodes them again to
// verify re-encoded data matches the originally encoded data.
func TestABIRoundtripEncodingAllTypes(t *testing.T) {
	// Create a value generator
	valueGenerator := NewRandomValueGenerator(&RandomValueGeneratorConfig{
		RandomArrayMinSize:  3,
		RandomArrayMaxSize:  10,
		RandomBytesMinSize:  5,
		RandomBytesMaxSize:  200,
		RandomStringMinSize: 5,
		RandomStringMaxSize: 200,
	}, rand.New(rand.NewSource(time.Now().UnixNano())))

	// Obtain our test ABI arguments
	args := getTestABIArguments()

	// Loop for each input argument
	for _, arg := range args {
		// Test each argument round trip serialization with different generated values (iterate a number of times).
		for i := 0; i < 10; i++ {
			// Generate a value for this argument
			value := GenerateAbiValue(valueGenerator, &arg.Type)

			// Encode the generated value for this argument
			encodedValue, err := encodeJSONArgument(&arg.Type, value)
			assert.NoError(t, err)

			// Decode the generated value
			decodedValue, err := decodeJSONArgument(&arg.Type, encodedValue, nil)
			assert.NoError(t, err)

			// Re-encode the generated value for this argument
			reencodedValue, err := encodeJSONArgument(&arg.Type, decodedValue)
			assert.NoError(t, err)

			// Compare the encoded and re-encoded values.
			matched := reflect.DeepEqual(encodedValue, reencodedValue)
			assert.True(t, matched, "round trip encoded->decoded->re-encoded ABI argument values did not match for '%v'.\nENCODED1: %v\nENCODED2: %v\n", arg.Name, encodedValue, reencodedValue)
		}
	}
}

// TestABIGenerationAndMutation runs tests to ABI value encoding works round-trip for argument values of all types.
// It generates values using a ValueGenerator, then encodes them, decodes them, and re-encodes them again to ensure
// re-encoded data matches the originally encoded data.
func TestABIGenerationAndMutation(t *testing.T) {
	// Create a value generator
	valueGenConfig := &MutatingValueGeneratorConfig{
		MinMutationRounds: 0,
		MaxMutationRounds: 1,
		RandomAddressBias: 0.5,
		RandomIntegerBias: 0.5,
		RandomStringBias:  0.5,
		RandomBytesBias:   0.5,
		RandomValueGeneratorConfig: &RandomValueGeneratorConfig{
			RandomArrayMinSize:  0,
			RandomArrayMaxSize:  100,
			RandomBytesMinSize:  0,
			RandomBytesMaxSize:  100,
			RandomStringMinSize: 0,
			RandomStringMaxSize: 100,
		},
	}
	valueGenerator := NewMutatingValueGenerator(valueGenConfig, NewValueSet(), rand.New(rand.NewSource(time.Now().UnixNano())))

	// Obtain our test ABI arguments
	args := getTestABIArguments()

	// Loop for each input argument
	for _, arg := range args {
		// Test each argument round trip serialization with different generated values (iterate a number of times).
		for i := 0; i < 10; i++ {
			// Generate a value for this argument
			value := GenerateAbiValue(valueGenerator, &arg.Type)

			// Mutate and ensure no error occurred.
			mutatedValue, err := MutateAbiValue(valueGenerator, &arg.Type, value)
			assert.NoError(t, err)

			// Verify the types of the value and mutated value are the same
			assert.EqualValues(t, reflect.ValueOf(value).Type().String(), reflect.ValueOf(mutatedValue).Type().String())
		}
	}
}
