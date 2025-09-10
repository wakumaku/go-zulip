package main

import (
	"fmt"
	"strings"
)

func main() {
	// Input data as provided in the problem statement
	inputData := `#	eventid	c
1	FINDERSALES_SavedSearchNotification	130
2	PCCM_inboxNewMessageSummaryNotification	23
3	PCCI_emailOtp	20
4	FINDERSALES_customerInquiryConfirmation	10
5	PCCK_TYD_milestoneUpdate	7
6	PCCI_pidWelcomeProspect	6
7	DPORDER_orderReceipt	5
8	ECOM_invoice_new	2
9	FINDERSALES_customerInquiryConfirmation_icc	2
10	PCCI_selfRegistrationExistingAccount	2
11	CONTENTINTEGRATION_eventBookingConfirmationReminder_BE_fr_FR	1
12	CONTENTINTEGRATION_eventBookingInvitation_GB_en_GB	1
13	PCCI_reminderInvitationNonPccVehicle	1
14	PCCI_reminderInvitationNonPccVehicleExistingAccount	1
15	PCCI_confirmChangedEmailaddress	1
16	CONTENTINTEGRATION_eventBookingInvitation_US_en_US	1
17	PCCI_addOwnerVehicle_ExistingAccount	1
18	CHARGINGSERVICE_invoice	1
19	PCCI_pidWelcomeOwner	1
20	PCCK_TYD_productionPicturesAvailable	1
21	PCCI_infoChangedEmailaddress	1
22	CONTENTINTEGRATION_eventBookingRegistrationReminder_US_en_US	1
23	PCCI_terminateVehicle_Owner	1
24	PCCK_sendDocumentsVerification_ExistingAccount	1
25	PCCI_lockedAccount	1
26	MSF_ticketing_paymentOrderInvoice	1`

	markdownTable := generateMarkdownTable(inputData)
	fmt.Println(markdownTable)
}

// generateMarkdownTable converts tab-delimited data to a markdown table,
// omitting the first column (row numbers) as requested
func generateMarkdownTable(input string) string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) == 0 {
		return ""
	}

	var result strings.Builder
	
	// Process header line - skip first column (row number)
	headerParts := strings.Split(lines[0], "\t")
	if len(headerParts) < 2 {
		return ""
	}
	
	// Write header (eventid and c columns)
	result.WriteString("| ")
	result.WriteString(headerParts[1]) // eventid
	result.WriteString(" | ")
	result.WriteString(headerParts[2]) // c
	result.WriteString(" |\n")
	
	// Write separator
	result.WriteString("|---|---|\n")
	
	// Process data lines - skip first column (row number)
	for i := 1; i < len(lines); i++ {
		parts := strings.Split(lines[i], "\t")
		if len(parts) < 3 {
			continue
		}
		
		result.WriteString("| ")
		result.WriteString(parts[1]) // eventid
		result.WriteString(" | ")
		result.WriteString(parts[2]) // c
		result.WriteString(" |\n")
	}
	
	return result.String()
}