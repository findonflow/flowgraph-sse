package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/r3labs/sse"
)

func main() {

	subscriptionQuery := ` subscription {
  latestTransaction {
    hash
    height
    index
    status
    keyIndex
    sequenceNumber
    gasLimit
    script
    arguments
    hasError
    error
    eventCount
    time
    payer {
      address
    }
    
    proposer {
      address
    }
    block {
      height
      time
    }
    referenceBlock {
      height
      time

    }
    contractInteractions {
      id
      locked
      deleted
      type
      address
      identifier
    }
    events(first: 50) {
      edges{
        node {
          fields
          type {
            fields{
              identifier
              type
            }
            id 
            name
          }
        }
      }
      pageInfo {
        hasNextPage
      }
    }
    tokenTransfers(
      first: 50
    ) {
      edges {
        node {
          type
          account {
            address
            domainNames {
              name
              provider
              fullName
            }
          }
          counterpartiesCount
          counterparty {
            address
            domainNames {
              name
              provider
              fullName
            }
          }
          amount {
            token {
              id
            }
            value
            usdValue
          }
        }
      }
      pageInfo{
        hasNextPage
      }
    }
    nftTransfers(
      first: 50
  
    ) {
      edges{
        node {
          from {
            address
            domainNames {
              name
              provider
              fullName
            }
          }
        }
      }
      pageInfo { 
        hasNextPage
      }
    }
  }
}`
	token := os.Getenv("FLOWGRAPH_KEY")
	if token == "" {
		panic("could not find flowgraph key")
	}
	url := fmt.Sprintf(`https://query.flowgraph.co/?query=%s&token=%s`, url.QueryEscape(subscriptionQuery), token)

	client := sse.NewClient(url)

	client.Subscribe("messages", func(msg *sse.Event) {
		// Got some data!
		var data interface{}
		err := json.Unmarshal(msg.Data, &data)
		if err != nil {
			fmt.Printf("[warn] %v", err)
		}
		b, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			fmt.Printf("[warn] %v", err)
		}
		os.Stdout.Write(b)

	})
}
