package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/0xbase-Corp/portfolio_svc/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//	@BasePath	/api/v1

// BitcoinController godoc
//
// @Summary      Fetch Bitcoin Wallet Information
// @Description  Retrieves information for a given Bitcoin address using the BTC.com API.
// @Tags         bitcoin
// @Accept       json
// @Produce      json
// @Param        btc-address  path      string  true  "Bitcoin Address"
// @Success      200          {object}  struct{ Data models.BitcoinAddressInfo "data"; ErrorCode int "error_code"; ErrNo int "err_no"; Message string "message"; Status string "status" }  "Returns wallet information including transaction history and balance"
// @Failure      400          {object}  struct{ Error string }  "Bad Request"
// @Failure      500          {object}  struct{ Error string }  "Internal Server Error"
// @Router       /portfolio/btc/:btc-address [get]

// BitcoinController handles requests for Bitcoin wallet information.
func BitcoinController(c *gin.Context, db *gorm.DB) {
	log.Println("BitcoinController invoked")
	// Extract the Bitcoin address from the request parameter.
	btcAddress := c.Param("btc-address")
	log.Println("for address:", btcAddress)

	// Prepare the BTC API request URL.
	url := "https://chain.api.btc.com/v3/address/" + btcAddress
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Create an HTTP client and execute the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	// Read the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("body ", body)
	var btcResponse models.BtcChainAPI
	if err := json.Unmarshal(body, &btcResponse); err != nil {
		log.Fatal(err)
	}

	fmt.Println("BTC PResponse ", btcResponse)
	//TODO: add the data into database

	c.JSON(http.StatusOK, gin.H{
		"data": btcResponse,
	})
}
