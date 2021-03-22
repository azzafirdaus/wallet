package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// HandleInitSession -> Initialize my account for wallet
func HandleInitSession(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Status: statusSuccess,
	}

	defer func() {
		json.NewEncoder(w).Encode(response)
	}()

	xid := r.FormValue("customer_xid")
	if xid == "" {
		response.Status = statusFail
		response.Data = ResponseInitAccount{
			Error: &ResponseInitAccountError{
				CustomerXID: []string{"Missing data for required field."},
			}}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sID, err := InitAccount(xid)
	if err != nil {
		response.Status = statusFail
		response.Data = ResponseInitAccount{
			Error: &ResponseInitAccountError{
				CustomerXID: []string{err.Error()},
			}}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Data = ResponseInitAccount{
		Token: sID,
	}
	w.WriteHeader(http.StatusCreated)
}

// HandleEnableWallet -> Enable my wallet
func HandleEnableWallet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := Response{
		Status: statusSuccess,
	}

	defer func() {
		json.NewEncoder(w).Encode(response)
	}()

	uID := r.FormValue("user_id")

	status, wallet, err := EnableWallet(uID)
	if err != nil {
		response.Status = statusFail
		response.Data = ResponseError{
			Error: err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !status {
		response.Status = statusFail
		response.Data = ResponseError{
			Error: "Already enabled",
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response.Data = ResponseWallet{
		Wallet: ResponseWalletDetail{
			ID:        wallet.ID,
			OwnedBy:   wallet.UserID,
			Status:    "enabled",
			EnabledAt: &wallet.EnableTime,
			Balance:   wallet.Balance,
		},
	}
	w.WriteHeader(http.StatusCreated)
}

// HandleViewBalance -> View my wallet balance
func HandleViewBalance(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := Response{
		Status: statusSuccess,
	}

	defer func() {
		json.NewEncoder(w).Encode(response)
	}()

	uID := r.FormValue("user_id")

	status, wallet, err := ViewBalance(uID)
	if err != nil {
		response.Status = statusFail
		response.Data = ResponseError{
			Error: err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !status {
		response.Status = statusFail
		response.Data = ResponseError{
			Error: "Disabled",
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response.Data = ResponseWallet{
		Wallet: ResponseWalletDetail{
			ID:        wallet.ID,
			OwnedBy:   wallet.UserID,
			Status:    "enabled",
			EnabledAt: &wallet.EnableTime,
			Balance:   wallet.Balance,
		},
	}
	w.WriteHeader(http.StatusOK)
}

// HandleDeposits -> Add virtual money to my wallet
func HandleDeposits(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := Response{
		Status: statusSuccess,
	}

	defer func() {
		json.NewEncoder(w).Encode(response)
	}()

	uID := r.FormValue("user_id")
	referenceID := r.FormValue("reference_id")
	amount, err := strconv.Atoi(r.FormValue("amount"))
	if err != nil {
		response.Status = statusFail
		response.Data = ResponseError{
			Error: "error read amount: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if referenceID == "" || amount < 0 {
		response.Status = statusFail
		response.Data = ResponseError{
			Error: "incorrect input",
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	status, tx, err := Deposit(uID, referenceID, amount)
	if err != nil {
		response.Status = statusFail
		response.Data = ResponseError{
			Error: err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !status {
		response.Status = statusFail
		response.Data = ResponseError{
			Error: "Wallet disabled",
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response.Data = ResponseDeposit{
		Deposit: ResponseDepositDetail{
			ID:          tx.ID,
			DepositedBy: uID,
			Status:      statusSuccess,
			DepositedAt: tx.CreateTime,
			Amount:      tx.Amount,
		},
	}
	w.WriteHeader(http.StatusCreated)
}

// HandleWithdrawal -> Use virtual money from my wallet
func HandleWithdrawal(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := Response{
		Status: statusSuccess,
	}

	defer func() {
		json.NewEncoder(w).Encode(response)
	}()

	uID := r.FormValue("user_id")
	referenceID := r.FormValue("reference_id")
	amount, err := strconv.Atoi(r.FormValue("amount"))
	if err != nil {
		response.Status = statusFail
		response.Data = ResponseError{
			Error: "error read amount: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if referenceID == "" || amount < 0 {
		response.Status = statusFail
		response.Data = ResponseError{
			Error: "incorrect input",
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	status, tx, err := Withdrawal(uID, referenceID, amount)
	if err != nil {
		response.Status = statusFail
		response.Data = ResponseError{
			Error: err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !status {
		response.Status = statusFail
		response.Data = ResponseError{
			Error: "Wallet disabled",
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response.Data = ResponseWithdrawal{
		Withdrawal: ResponseWithdrawalDetail{
			ID:          tx.ID,
			WithdrawnBy: uID,
			Status:      statusSuccess,
			WithdrawnAt: tx.CreateTime,
			Amount:      tx.Amount,
		},
	}
	w.WriteHeader(http.StatusCreated)
}

// HandleDisableWallet -> Disable my wallet
func HandleDisableWallet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := Response{
		Status: statusSuccess,
	}

	defer func() {
		json.NewEncoder(w).Encode(response)
	}()

	uID := r.FormValue("user_id")

	status, wallet, err := DisableWallet(uID)
	if err != nil {
		response.Status = statusFail
		response.Data = ResponseError{
			Error: err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !status {
		response.Status = statusFail
		response.Data = ResponseError{
			Error: "Already disabled",
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response.Data = ResponseWallet{
		Wallet: ResponseWalletDetail{
			ID:         wallet.ID,
			OwnedBy:    wallet.UserID,
			Status:     "disabled",
			DisabledAt: &wallet.EnableTime,
			Balance:    wallet.Balance,
		},
	}
	w.WriteHeader(http.StatusCreated)
}
