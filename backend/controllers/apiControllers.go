/*
* Package controllers file for API
* @author: Ayan Banerjee
* @Organization: Math & Cody
 */
package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"ethential/dapp/utils"

	"ethential/dapp/types"

	"github.com/gin-gonic/gin"
)

var (
	WALLET_URI = utils.GoDotEnvVariable("WALLET_URI")
	MNGR_ADDR  = utils.GoDotEnvVariable("MNGR_CONTRACT_ADDR")
)

// GenTokenController -> to generate new token with the auth service
func GenTokenController(c *gin.Context) {
	clientID := c.Param("clientID")
	token, err := utils.GenToken(clientID)
	if err != nil {
		c.JSON(500, gin.H{
			"Error":        "Error Generating Token",
			"ErrorDetails": err.Error(),
		})
		return
	}
	c.JSON(200, token)
	return
}

// TransferTokenController -> To create unsigned tx controller and store it to redis with txHash as key and unsigned tx as param
func TransferTokenController(c *gin.Context) {
	h := types.AuthHeaderType{}
	if err := c.ShouldBindHeader(&h); err != nil {
		c.JSON(502, gin.H{
			"Error":        "Error Parsing Header",
			"ErrorDetails": err.Error(),
		})
		return
	}
	if strings.Contains(h.Token, "Bearer ") {
		token := strings.Replace(h.Token, "Bearer ", "", -1)
		isValidToken := utils.VerifyToken(token)
		if isValidToken {
			params := []types.Params{}
			reqParams := types.ReqBody{}
			c.BindJSON(&reqParams)
			p := types.Params{
				InternalType: "address",
				Name:         "recipient",
				Type:         "address",
				Value:        reqParams.ToAddress,
			}
			params = append(params, p)
			p = types.Params{
				InternalType: "uint256",
				Name:         "amount",
				Type:         "uint256",
				Value:        reqParams.TokenAmount,
			}
			params = append(params, p)
			txParams := types.TxParams{
				From:   reqParams.From,
				Params: params,
				Value:  0,
			}
			unsignedTx, err := utils.CreateTransferTokenTx(token, txParams, 5)
			if err != nil {
				c.JSON(500, gin.H{
					"Error": err.Error(),
				})
				return
			}
			// NOTE - Make sure the final Transaction is added to the unsignedTx variable
			postData := map[string]string{
				"tx": unsignedTx,
			}
			postDataByte, err := json.Marshal(postData)
			if err != nil {
				c.JSON(500, gin.H{
					"Error": "Error Marshalling POST JSON Data!",
				})
				return
			}
			resp, err := http.Post(WALLET_URI, "application/json", bytes.NewBuffer(postDataByte))
			if err != nil {
				c.JSON(500, gin.H{
					"Error": err.Error(),
				})
			}
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			if resp.StatusCode == 200 {
				txHash := string(body)
				// Sending TxHash as of now.
				c.JSON(200, txHash)
				return
			}
			c.String(resp.StatusCode, string(body))
			return
		} else {
			c.JSON(403, gin.H{
				"Error": "Invalid AuthToken!",
			})
			return
		}
	} else {
		c.JSON(403, gin.H{
			"Error": "No AuthToken Supplied!",
		})
		return
	}
}

// TokenBalanceController -> Gets token balance for given public address
func TokenBalanceController(c *gin.Context) {
	h := types.AuthHeaderType{}
	if err := c.ShouldBindHeader(&h); err != nil {
		c.JSON(502, gin.H{
			"Error":        "Error Parsing Header",
			"ErrorDetails": err.Error(),
		})
		return
	}
	if strings.Contains(h.Token, "Bearer ") {
		token := strings.Replace(h.Token, "Bearer ", "", -1)
		isValidToken := utils.VerifyToken(token)
		if isValidToken {
			params := []types.Params{}
			reqParams := types.ReqBody{}
			c.BindJSON(&reqParams)
			p := types.Params{
				InternalType: "address",
				Name:         "account",
				Type:         "address",
				Value:        reqParams.Pubkey,
			}
			fmt.Println(p)
			params = append(params, p)
			txParams := types.TxParams{
				From:   reqParams.Pubkey,
				Params: params,
				Value:  0,
			}
			balance, err := utils.GetTokenBalance(token, txParams, 5)
			if err != nil {
				c.JSON(500, gin.H{
					"Error": err.Error(),
				})
				return
			}
			c.JSON(200, balance)
		}
	}
}

// SwapTokenController -> Swap a token for wrapped ether
func SwapTokenController(c *gin.Context) {
	h := types.AuthHeaderType{}
	if err := c.ShouldBindHeader(&h); err != nil {
		c.JSON(502, gin.H{
			"Error":        "Error Parsing Header",
			"ErrorDetails": err.Error(),
		})
		return
	}
	if strings.Contains(h.Token, "Bearer ") {
		token := strings.Replace(h.Token, "Bearer ", "", -1)
		isValidToken := utils.VerifyToken(token)
		if isValidToken {
			params := []types.Params{}
			reqParams := types.ReqBody{}
			c.BindJSON(&reqParams)
			p := types.Params{
				InternalType: "uint256",
				Name:         "amount",
				Type:         "uint256",
				Value:        reqParams.TokenAmount,
			}
			fmt.Println(p)
			params = append(params, p)
			txParams := types.TxParams{
				From:   reqParams.Pubkey,
				Params: params,
				Value:  0,
			}
			unsignedTx, err := utils.CreateSwapTokenTx(token, txParams, 5)
			if err != nil {
				c.JSON(500, gin.H{
					"Error": err.Error(),
				})
				return
			}
			postData := map[string]string{
				"tx": unsignedTx,
			}
			postDataByte, err := json.Marshal(postData)
			if err != nil {
				c.JSON(500, gin.H{
					"Error": "Error Marshalling POST JSON Data!",
				})
				return
			}
			resp, err := http.Post(WALLET_URI, "application/json", bytes.NewBuffer(postDataByte))
			if err != nil {
				c.JSON(500, gin.H{
					"Error": err.Error(),
				})
			}
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			if resp.StatusCode == 200 {
				txHash := string(body)
				// Sending TxHash as of now.
				c.JSON(200, txHash)
				return
			}
			c.String(resp.StatusCode, string(body))
			return
		}
	}
}

func ApproveController(c *gin.Context) {
	h := types.AuthHeaderType{}
	if err := c.ShouldBindHeader(&h); err != nil {
		c.JSON(502, gin.H{
			"Error":        "Error Parsing Header",
			"ErrorDetails": err.Error(),
		})
		return
	}
	if strings.Contains(h.Token, "Bearer ") {
		token := strings.Replace(h.Token, "Bearer ", "", -1)
		isValidToken := utils.VerifyToken(token)
		if isValidToken {
			params := []types.Params{}
			reqParams := types.ReqBody{}
			c.BindJSON(&reqParams)
			fmt.Println(reqParams)
			p := types.Params{
				InternalType: "address",
				Name:         "spender",
				Type:         "address",
				Value:        MNGR_ADDR,
			}
			params = append(params, p)
			p = types.Params{
				InternalType: "uint256",
				Name:         "amount",
				Type:         "uint256",
				Value:        reqParams.TokenAmount,
			}
			params = append(params, p)
			fmt.Println(params)
			txParams := types.TxParams{
				From:   reqParams.Pubkey,
				Params: params,
				Value:  0,
			}
			unsignedTx, err := utils.CreateApproveTx(token, txParams, 5)
			if err != nil {
				fmt.Println(err)
				c.JSON(500, gin.H{
					"Error": err.Error(),
				})
				return
			}
			postData := map[string]string{
				"tx": unsignedTx,
			}
			postDataByte, err := json.Marshal(postData)
			if err != nil {
				c.JSON(500, gin.H{
					"Error": "Error Marshalling POST JSON Data!",
				})
				return
			}
			resp, err := http.Post(WALLET_URI, "application/json", bytes.NewBuffer(postDataByte))
			if err != nil {
				c.JSON(500, gin.H{
					"Error": err.Error(),
				})
			}
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			if resp.StatusCode == 200 {
				txHash := string(body)
				// Sending TxHash as of now.
				c.JSON(200, txHash)
				return
			}
			c.String(resp.StatusCode, string(body))
			return
		}
	}
}
