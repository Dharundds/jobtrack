package classify

import "testing"

func TestExtractCompany_greenhouseSharedSender(t *testing.T) {
	from := "no-reply@us.greenhouse-mail.io"
	rubrik := ExtractCompany(from, "Thank you for applying to Rubrik", "Thanks for applying to Rubrik")
	if rubrik != "Rubrik" {
		t.Fatalf("Rubrik company = %q", rubrik)
	}
	prodigal := ExtractCompany(from, "Thank you for applying to Prodigal", "Thanks for applying to Prodigal")
	if prodigal != "Prodigal" {
		t.Fatalf("Prodigal company = %q", prodigal)
	}
	if SimilarCompanyRole(rubrik, "", prodigal, "") {
		t.Fatal("Rubrik and Prodigal should not match")
	}
}

func TestSimilarCompanyRole_ATSBrand(t *testing.T) {
	if SimilarCompanyRole("Greenhouse", "", "Greenhouse", "") {
		t.Fatal("should not merge on ATS brand alone")
	}
}
