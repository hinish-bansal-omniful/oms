package parse_csv

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"oms-service/domain"
	"oms-service/intersvc"

	"github.com/omniful/go_commons/csv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ParseCSV reads a CSV file and converts it into a slice of domain.Order objects
func ParseCSV(filePath string) ([]*domain.Order, error) {
	fmt.Println("ParseCSV function called successfully!")
	fmt.Println("Opening file: ", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Map to group items by order_id
	orderGroups := make(map[string]*domain.Order)

	// Initialize CSV reader
	CSV, err := csv.NewCommonCSV(
		csv.WithBatchSize(100),
		csv.WithSource(csv.Local),
		csv.WithLocalFileInfo(filePath),
		csv.WithHeaderSanitizers(csv.SanitizeAsterisks, csv.SanitizeToLower),
		csv.WithDataRowSanitizers(csv.SanitizeSpace, csv.SanitizeToLower),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize CSV reader: %v", err)
	}

	err = CSV.InitializeReader(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize CSV reader: %v", err)
	}

	// Process the records and group them by order_id
	for !CSV.IsEOF() {
		records, err := CSV.ReadNextBatch()
		if err != nil {
			log.Println("Error reading batch:", err)
			continue
		}

		fmt.Println("Processing records:", records)

		for _, record := range records {
			if len(record) < 5 {
				log.Println("Skipping invalid row (not enough columns):", record)
				continue
			}

			orderID := record[0]  // order_id
			skuID := record[1]    // sku_id
			quantity := record[2] // quantity
			sellerID := record[3] // seller_id
			hubID := record[4]    // hub_id

			// Convert quantity to integer
			intQuantity, err := strconv.Atoi(quantity)
			if err != nil {
				log.Println("Skipping invalid quantity:", quantity, err)
				continue
			}

			// Simulated price assignment (could be fetched from DB)
			price := 100.0 // Example static price; replace with dynamic logic if needed

			// Check if the order already exists
			order, exists := orderGroups[orderID]
			if !exists {
				now := time.Now()
				order = &domain.Order{
					ID:          primitive.NewObjectID(), // Generate a new ObjectID
					OrderID:     orderID,
					TenantID:    "111", // Example value, replace as needed
					SellerID:    sellerID,
					HubID:       hubID,
					CustomerID:  sellerID, // Assuming customer_id is same as seller_id
					Items:       []domain.OrderItem{},
					TotalAmount: 0,
					Status:      "on_hold",
					CreatedAt:   now,
					UpdatedAt:   now,
				}
				orderGroups[orderID] = order
			}

			// Create OrderItem
			orderItem := domain.OrderItem{
				ID:       primitive.NewObjectID(),
				SKUID:    skuID,
				Quantity: intQuantity,
				Price:    price,
			}
			order.Items = append(order.Items, orderItem)
		}
	}

	// Convert orderGroups map to slice
	var orders []*domain.Order
	for _, order := range orderGroups {
		orders = append(orders, order)
	}

	fmt.Println("Final orders:")
	for _, order := range orders {
		fmt.Printf("Order No: %s, Total Items: %d\n", order.OrderID, len(order.Items))
		go intersvc.ValidateOrders(order)
	}

	return orders, nil
}

// package parse_csv

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"oms-service/intersvc"
// 	"oms-service/domain"
// 	"os"
// 	"strconv"
// 	"time"

// 	"github.com/omniful/go_commons/csv"
// )

// func ParseCSV(filePath string) ([]*domain.Order, error) {
// 	fmt.Println("Parse CSV function called successfull!")
// 	fmt.Println("This is the file beign opened: ", filePath)
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		fmt.Println("Error in opening file path.")
// 	}
// 	defer file.Close()

// 	// Map to group items by order_no
// 	orderGroups := make(map[string]*domain.Order)

// 	// Initialize the CSV reader (based on your previous implementation)
// 	CSV, err := csv.NewCommonCSV(
// 		csv.WithBatchSize(100),
// 		csv.WithSource(csv.Local),
// 		csv.WithLocalFileInfo(filePath),
// 		csv.WithHeaderSanitizers(csv.SanitizeAsterisks, csv.SanitizeToLower),
// 		csv.WithDataRowSanitizers(csv.SanitizeSpace, csv.SanitizeToLower),
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to initialize CSV reader: %v", err)
// 	}

// 	err = CSV.InitializeReader(context.TODO())
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to initialize CSV reader: %v", err)
// 	}

// 	// Process the records and group them by order_no and customer_name
// 	for !CSV.IsEOF() {
// 		var records csv.Records
// 		records, err := CSV.ReadNextBatch()
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		fmt.Println("Processing records:")
// 		fmt.Println(records)
// 		for _, record := range records {
// 			orderID := record[0]  // order_id
// 			skuID := record[1]    // sku_id
// 			quantity := record[2] // quantity
// 			sellerID := record[3] // seller_id
// 			hubID := record[4]    // hub_id

// 			// Convert quantity to integer
// 			IntQuantity, err := strconv.Atoi(quantity)
// 			if err != nil {
// 				return nil, fmt.Errorf("invalid quantity %s: %v", quantity, err)
// 			}

// 			// Check if the  group forderor this order_id already exists
// 			orderKey := orderID
// 			order, exists := orderGroups[orderKey]
// 			if !exists {
// 				// If order doesn't exist, create a new order
// 				now := time.Now()
// 				order = &domain.Order{
// 					// SellerID:     sellerID,
// 					// HubID:        hubID,
// 					ID:              orderID,
// 					CustomerID:      sellerID,
// 					CreatedAt:       now,
// 					Currency:        "INR",
// 					TotalAmount:     0,
// 					TransactionID:   "SAMPLE_TXN_ID",
// 					ModeOfPayment:   "PAYPAL",
// 					Status:          "on_hold",
// 					BillingAddress:  "sample address",
// 					ShippingAddress: "sample address",
// 					InvoiceID:       999,
// 					TenantID:        111,
// 					OrderItems:      []domain.OrderItem{}, // Start with an empty slice of items
// 				}
// 				// Add the new order to the map
// 				orderGroups[orderKey] = order
// 			}

// 			// Create a new OrderItem and append it to the order's OrderItems
// 			orderItem := domain.OrderItem{
// 				OrderID:         orderID,
// 				SKUID:           skuID,
// 				Quantity: IntQuantity,
// 				Price: price
// 			}
// 			order.OrderItems = append(order.OrderItems, orderItem)
// 		}
// 	}

// 	var orders []*domain.Order
// 	for _, order := range orderGroups {
// 		orders = append(orders, order)
// 	}

// 	fmt.Println("Final orders:")
// 	for _, order := range orders {
// 		fmt.Printf("Order No: %s, Total Items: %d\n", order.ID, len(order.OrderItems))
// 		go intersvc.ValidateOrders(order)
// 	}

// 	return orders, nil

// }
