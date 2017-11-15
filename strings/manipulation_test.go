package strings

import (
	"testing"
)

func TestToASCII(t *testing.T) {
	var input, cleaned string

	input = "你好，世界"
	cleaned = ToASCII(input)
	if "Ni_Hao_Shi_Jie" != cleaned {
		t.Fatal("ToASCII failed to clean Chinese UTF-8.")
	}

	input = "Nǐ hǎo, shìjiè"
	cleaned = ToASCII(input)
	if "Ni_hao_shijie" != cleaned {
		t.Fatal("ToASCII failed to clean Chinese spoken UTF-8.")
	}

	input = "こんにちは世界"
	cleaned = ToASCII(input)
	if "konnichihaShi_Jie" != cleaned {
		t.Fatal("ToASCII failed to clean Japanese UTF-8.")
	}

	input = "Kon'nichiwa sekai"
	cleaned = ToASCII(input)
	if "Konnichiwa_sekai" != cleaned {
		t.Fatal("ToASCII failed to clean Japanese spoken UTF-8.")
	}

	input = "Привет мир"
	cleaned = ToASCII(input)
	if "Privet_mir" != cleaned {
		t.Fatal("ToASCII failed to clean Russian UTF-8.")
	}

	input = "totally_legit"
	cleaned = ToASCII(input)
	if input != cleaned {
		t.Fatal("ToASCII improperly cleaned a totally legit name.")
	}

	input = "  _ too  _() much (*&%$$()*^)*$ __ extra  _  _	"
	cleaned = ToASCII(input)
	if "too_much_extra" != cleaned {
		t.Fatal("ToASCII improperly cleaned a name with too much extra stuff.")
	}
}

func TestSanitizeFolderPath(t *testing.T) {
	var input, cleaned string

	//Sanitize folder/file tests....
	input = "/Привет sdflkj_-2345/myfilename&:*_.img"
	cleaned = SanitizeFolderPath(input)
	if "/Privet_sdflkj_-2345/myfilename_.img" != cleaned {
		t.Fatalf("SanitizeFolderPath improperly cleaned a name with too much extra stuff: %s", cleaned)
	}

	input = "Привет sdflkj_-2345/my filename&*_.img"
	cleaned = SanitizeFolderPath(input)
	if "Privet_sdflkj_-2345/my_filename_.img" != cleaned {
		t.Fatalf("SanitizeFolderPath improperly cleaned a name with too much extra stuff: %s", cleaned)
	}
}

func TestDigitsOnly(t *testing.T) {
	var s string

	s = "1 2 3 4 5 6 7 8 9 0"
	s = DigitsOnly(s)
	if s != "1234567890" {
		t.Fatalf("DigitsOnly failed (01). Got:->%s<-", s)
	}

	s = "!@#$%^&*()QWERTYUIOPASDFGHJKLZXCVBNM{}|_+:"
	s = DigitsOnly(s)
	if s != "" {
		t.Fatalf("DigitsOnly failed (02). Got:->%s<-", s)
	}
}
