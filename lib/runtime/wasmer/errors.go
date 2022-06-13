package wasmer

import (
	"errors"
	"fmt"
	"github.com/ChainSafe/gossamer/lib/transaction"
	"github.com/ChainSafe/gossamer/pkg/scale"
)

var (
	errInvalidResult = errors.New("invalid error value")
	errInvalidType   = errors.New("invalid validity type")
)

func newUnknownError(data scale.VaryingDataTypeValue) error {
	return fmt.Errorf("unknown error: %d", data)
}

// UnmarshalError occurs when unmarshalling fails
type UnmarshalError struct {
	msg string
}

func (e UnmarshalError) Error() string {
	return fmt.Sprintf("unmarshal error: %s", e.msg)
}

// A TransactionValidityError is possible errors while checking the validity of a transaction
type TransactionValidityError struct {
	msg error // description of error
}

func (e TransactionValidityError) Error() string {
	return fmt.Sprintf("transaction validity error: %s", e.msg)
}

// Other Some error occurred
type Other string

// Index Returns VDT index
func (err Other) Index() uint { return 0 }

// CannotLookup Failed to lookup some data
type CannotLookup struct{}

// Index Returns VDT index
func (err CannotLookup) Index() uint { return 1 }

// BadOrigin A bad origin
type BadOrigin struct{}

// Index Returns VDT index
func (err BadOrigin) Index() uint { return 2 }

// Module A custom error in a module
type Module struct {
	Idx     uint8
	Err     uint8
	Message *string
}

// Index Returns VDT index
func (err Module) Index() uint { return 3 }

func (err Module) string() string {
	return fmt.Sprintf("index: %d code: %d message: %x", err.Idx, err.Err, *err.Message)
}

// ValidityCannotLookup Could not lookup some information that is required to validate the transaction
type ValidityCannotLookup struct{}

// Index Returns VDT index
func (err ValidityCannotLookup) Index() uint { return 0 }

// NoUnsignedValidator No validator found for the given unsigned transaction
type NoUnsignedValidator struct{}

// Index Returns VDT index
func (err NoUnsignedValidator) Index() uint { return 1 }

// UnknownCustom Any other custom unknown validity that is not covered
type UnknownCustom uint8

// Index Returns VDT index
func (err UnknownCustom) Index() uint { return 2 }

// Call The call of the transaction is not expected
type Call struct{}

// Index Returns VDT index
func (err Call) Index() uint { return 0 }

// Payment General error to do with the inability to pay some fees (e.g. account balance too low)
type Payment struct{}

// Index Returns VDT index
func (err Payment) Index() uint { return 1 }

// Future General error to do with the transaction not yet being valid (e.g. nonce too high)
type Future struct{}

// Index Returns VDT index
func (err Future) Index() uint { return 2 }

// Stale General error to do with the transaction being outdated (e.g. nonce too low)
type Stale struct{}

// Index Returns VDT index
func (err Stale) Index() uint { return 3 }

// BadProof General error to do with the transaction’s proofs (e.g. signature)
type BadProof struct{}

// Index Returns VDT index
func (err BadProof) Index() uint { return 4 }

// AncientBirthBlock The transaction birth block is ancient
type AncientBirthBlock struct{}

// Index Returns VDT index
func (err AncientBirthBlock) Index() uint { return 5 }

// ExhaustsResources The transaction would exhaust the resources of current block
type ExhaustsResources struct{}

// Index Returns VDT index
func (err ExhaustsResources) Index() uint { return 6 }

// InvalidCustom Any other custom invalid validity that is not covered
type InvalidCustom uint8

// Index Returns VDT index
func (err InvalidCustom) Index() uint { return 7 }

// BadMandatory An extrinsic with a Mandatory dispatch resulted in Error
type BadMandatory struct{}

// Index Returns VDT index
func (err BadMandatory) Index() uint { return 8 }

// MandatoryDispatch A transaction with a mandatory dispatch
type MandatoryDispatch struct{}

// Index Returns VDT index
func (err MandatoryDispatch) Index() uint { return 9 }

// BadSigner The sending address is disabled or known to be invalid.
type BadSigner struct{}

// Index Returns VDT index
func (err BadSigner) Index() uint { return 10 }

// API Errors
type FailedToDecodeReturnValue string

func (err FailedToDecodeReturnValue) Index() uint { return 0 }

type FailedToConvertReturnValue string

func (err FailedToConvertReturnValue) Index() uint { return 1 }

type FailedToConvertParameter string

func (err FailedToConvertParameter) Index() uint { return 2 }

type Application string

func (err Application) Index() uint { return 3 }

var (
	errUnexpectedTxCall          = errors.New("call of the transaction is not expected")
	errInvalidPayment            = errors.New("invalid payment")
	errInvalidTransaction        = errors.New("invalid transaction")
	errOutdatedTransaction       = errors.New("outdated transaction")
	errBadProof                  = errors.New("bad proof")
	errAncientBirthBlock         = errors.New("ancient birth block")
	errExhaustsResources         = errors.New("exhausts resources")
	errMandatoryDispatchError    = errors.New("mandatory dispatch error")
	errInvalidMandatoryDispatch  = errors.New("invalid mandatory dispatch")
	errLookupFailed              = errors.New("lookup failed")
	errValidatorNotFound         = errors.New("validator not found")
	errBadSigner                 = errors.New("invalid signing address")
	invalidCustom                InvalidCustom
	unknownCustom                UnknownCustom
	errFailedToDecodeReturnValue FailedToDecodeReturnValue
)

func determineErrType(vdt scale.VaryingDataType) error {
	switch val := vdt.Value().(type) {
	// InvalidTransaction Error
	case Call:
		return &TransactionValidityError{errUnexpectedTxCall}
	case Payment:
		return &TransactionValidityError{errInvalidPayment}
	case Future:
		return &TransactionValidityError{errInvalidTransaction}
	case Stale:
		return &TransactionValidityError{errOutdatedTransaction}
	case BadProof:
		return &TransactionValidityError{errBadProof}
	case AncientBirthBlock:
		return &TransactionValidityError{errAncientBirthBlock}
	case ExhaustsResources:
		return &TransactionValidityError{errExhaustsResources}
	case InvalidCustom:
		return &TransactionValidityError{newUnknownError(val)}
	case BadMandatory:
		return &TransactionValidityError{errMandatoryDispatchError}
	case MandatoryDispatch:
		return &TransactionValidityError{errInvalidMandatoryDispatch}
	case BadSigner:
		return &TransactionValidityError{errBadSigner}

		// UnknownTransaction Error
	case ValidityCannotLookup:
		return &TransactionValidityError{errLookupFailed}
	case NoUnsignedValidator:
		return &TransactionValidityError{errValidatorNotFound}
	case UnknownCustom:
		return &TransactionValidityError{newUnknownError(val)}
	}

	fmt.Println("Why am i hitting this??")
	return errInvalidResult
}

func decodeValidity(res []byte) (*transaction.Validity, error) {
	invalid := scale.MustNewVaryingDataType(Call{}, Payment{}, Future{}, Stale{}, BadProof{}, AncientBirthBlock{},
		ExhaustsResources{}, invalidCustom, BadMandatory{}, MandatoryDispatch{})
	unknown := scale.MustNewVaryingDataType(ValidityCannotLookup{}, NoUnsignedValidator{}, unknownCustom)
	apiErr := scale.MustNewVaryingDataType(errFailedToDecodeReturnValue)

	//Result<TransactionValidityResult, APIError>
	//TransactionValidityResult<TransactionValidity, TransactionValidityError

	validTxn := transaction.Validity{}
	txnValidityErrResult := scale.NewResult(invalid, unknown)
	txnValidityResult := scale.NewResult(validTxn, txnValidityErrResult)

	result := scale.NewResult(txnValidityResult, apiErr)

	err := scale.Unmarshal(res, &result)
	if err != nil {
		return nil, &UnmarshalError{err.Error()}
	}

	ok, err := result.Unwrap()
	if err != nil {
		fmt.Println("api err case")
		//TODO implement this
		// APIError
		switch err := err.(type) {
		default:
			fmt.Println(err)
			fmt.Println("d")
			return nil, errInvalidResult
		}
	} else {
		// TxnValidity
		fmt.Println("txn Validity case")
		switch o := ok.(type) {
		case scale.Result:
			// TxnValidityErr Result
			txnValidityRes, err := o.Unwrap()
			if err != nil {
				fmt.Println("txn validity error case")
				switch errType := err.(type) {

				// Err wrapping result
				case scale.WrappedErr:
					fmt.Println("in wrapped error")
					errResult := errType.Err.(scale.Result)
					ok, err = errResult.Unwrap()
					if err != nil {
						fmt.Println("unknown case")
						switch err := err.(type) {
						case scale.WrappedErr:
							return nil, determineErrType(err.Err.(scale.VaryingDataType))
						default:
							return nil, errInvalidResult
						}
					} else {
						fmt.Println("invalid case")
						return nil, determineErrType(ok.(scale.VaryingDataType))
					}
				default:
					fmt.Println("b")
					return nil, errInvalidResult
				}
			} else {
				switch validity := txnValidityRes.(type) {
				case transaction.Validity:
					return &validity, nil
				default:
					return nil, errInvalidType
				}
			}
		default:
			fmt.Println("a")
			return nil, errInvalidResult
		}
	}
}
