package apdu

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mniak/tlv"
	"github.com/stretchr/testify/assert"
)

func TestEMVProprietaryTemplate_Merge(t *testing.T) {
	getRandomTemplate := func() EMVProprietaryTemplate {
		var result EMVProprietaryTemplate
		gofakeit.Struct(&result)
		// Some fields are ignored for now because I don't know how would be the
		// correct way to merge its values
		result.ApplicationTemplates = nil

		// Ignore rawtlv because it must really be ignored
		result.RawTLV = nil
		return result
	}

	t.Run("Non-empty merge with empty", func(t *testing.T) {
		a := getRandomTemplate()
		var b_empty EMVProprietaryTemplate

		merged := a.Merge(b_empty)
		assert.Equal(t, a, merged)
	})
	t.Run("Non-empty merge with non-empty", func(t *testing.T) {
		a := getRandomTemplate()
		b := getRandomTemplate()

		merged := a.Merge(b)
		assert.Equal(t, a, merged)
	})
	t.Run("Empty merge with non-empty", func(t *testing.T) {
		var a_empty EMVProprietaryTemplate
		b := getRandomTemplate()

		merged := a_empty.Merge(b)
		assert.Equal(t, b, merged)
	})

	t.Run("Test field by field", func(t *testing.T) {
		testCases := []struct {
			field       string
			emptyValue  func() any
			randomValue func() any
			get         func(et *EMVProprietaryTemplate) any
			set         func(et *EMVProprietaryTemplate, val any)
		}{
			{
				field:       "Track2EquivalentData",
				emptyValue:  func() any { return "" },
				randomValue: func() any { return gofakeit.SentenceSimple() },
				get:         func(et *EMVProprietaryTemplate) any { return et.Track2EquivalentData },
				set:         func(et *EMVProprietaryTemplate, val any) { et.Track2EquivalentData = val.([]byte) },
			},
			{
				field:       "CardholderName",
				emptyValue:  func() any { return "" },
				randomValue: func() any { return gofakeit.SentenceSimple() },
				get:         func(et *EMVProprietaryTemplate) any { return et.CardholderName },
				set:         func(et *EMVProprietaryTemplate, val any) { et.CardholderName = val.(string) },
			},
			{
				field:       "PAN",
				emptyValue:  func() any { return "" },
				randomValue: func() any { return gofakeit.SentenceSimple() },
				get:         func(et *EMVProprietaryTemplate) any { return et.PAN },
				set:         func(et *EMVProprietaryTemplate, val any) { et.PAN = val.(string) },
			},
			{
				field:       "PANSequenceNumber",
				emptyValue:  func() any { return "" },
				randomValue: func() any { return gofakeit.SentenceSimple() },
				get:         func(et *EMVProprietaryTemplate) any { return et.PANSequenceNumber },
				set:         func(et *EMVProprietaryTemplate, val any) { et.PANSequenceNumber = val.(string) },
			},
			{
				field:       "ExpirationDate",
				emptyValue:  func() any { return "" },
				randomValue: func() any { return gofakeit.SentenceSimple() },
				get:         func(et *EMVProprietaryTemplate) any { return et.ExpirationDate },
				set:         func(et *EMVProprietaryTemplate, val any) { et.ExpirationDate = val.(string) },
			},
			{
				field:       "UsageControl",
				emptyValue:  func() any { return "" },
				randomValue: func() any { return gofakeit.SentenceSimple() },
				get:         func(et *EMVProprietaryTemplate) any { return et.UsageControl },
				set:         func(et *EMVProprietaryTemplate, val any) { et.UsageControl = val.(string) },
			},
			{
				field:       "IssuerCountryCode",
				emptyValue:  func() any { return "" },
				randomValue: func() any { return gofakeit.SentenceSimple() },
				get:         func(et *EMVProprietaryTemplate) any { return et.IssuerCountryCode },
				set:         func(et *EMVProprietaryTemplate, val any) { et.IssuerCountryCode = val.(string) },
			},
			{
				field:       "EffectiveDate",
				emptyValue:  func() any { return "" },
				randomValue: func() any { return gofakeit.SentenceSimple() },
				get:         func(et *EMVProprietaryTemplate) any { return et.EffectiveDate },
				set:         func(et *EMVProprietaryTemplate, val any) { et.EffectiveDate = val.(string) },
			},
			{
				field:       "IssuerActionCodeDenial",
				emptyValue:  func() any { return "" },
				randomValue: func() any { return gofakeit.SentenceSimple() },
				get:         func(et *EMVProprietaryTemplate) any { return et.IssuerActionCodeDenial },
				set:         func(et *EMVProprietaryTemplate, val any) { et.IssuerActionCodeDenial = val.(string) },
			},
			{
				field:       "IssuerActionCodeOnline",
				emptyValue:  func() any { return "" },
				randomValue: func() any { return gofakeit.SentenceSimple() },
				get:         func(et *EMVProprietaryTemplate) any { return et.IssuerActionCodeOnline },
				set:         func(et *EMVProprietaryTemplate, val any) { et.IssuerActionCodeOnline = val.(string) },
			},
			{
				field:       "IssuerActionCodeDefault",
				emptyValue:  func() any { return "" },
				randomValue: func() any { return gofakeit.SentenceSimple() },
				get:         func(et *EMVProprietaryTemplate) any { return et.IssuerActionCodeDefault },
				set:         func(et *EMVProprietaryTemplate, val any) { et.IssuerActionCodeDefault = val.(string) },
			},
			{
				field:       "CAPublicKeyIndex1",
				emptyValue:  func() any { return "" },
				randomValue: func() any { return gofakeit.SentenceSimple() },
				get:         func(et *EMVProprietaryTemplate) any { return et.CAPublicKeyIndex1 },
				set:         func(et *EMVProprietaryTemplate, val any) { et.CAPublicKeyIndex1 = val.(string) },
			},
			{
				field:       "IssuerPublicKeyExponent",
				emptyValue:  func() any { return "" },
				randomValue: func() any { return gofakeit.SentenceSimple() },
				get:         func(et *EMVProprietaryTemplate) any { return et.IssuerPublicKeyExponent },
				set:         func(et *EMVProprietaryTemplate, val any) { et.IssuerPublicKeyExponent = val.(string) },
			},
			{
				field:       "IssuerPublicKeyCertificate",
				emptyValue:  func() any { return "" },
				randomValue: func() any { return gofakeit.SentenceSimple() },
				get:         func(et *EMVProprietaryTemplate) any { return et.IssuerPublicKeyCertificate },
				set:         func(et *EMVProprietaryTemplate, val any) { et.IssuerPublicKeyCertificate = val.(string) },
			},
			{
				field:       "CurrencyCode",
				emptyValue:  func() any { return "" },
				randomValue: func() any { return gofakeit.SentenceSimple() },
				get:         func(et *EMVProprietaryTemplate) any { return et.CurrencyCode },
				set:         func(et *EMVProprietaryTemplate, val any) { et.CurrencyCode = val.(string) },
			},
			{
				field:       "CurrencyExponent",
				emptyValue:  func() any { return "" },
				randomValue: func() any { return gofakeit.SentenceSimple() },
				get:         func(et *EMVProprietaryTemplate) any { return et.CurrencyExponent },
				set:         func(et *EMVProprietaryTemplate, val any) { et.CurrencyExponent = val.(string) },
			},
			{
				field:       "CDOL1",
				emptyValue:  func() any { return make([]byte, 0) },
				randomValue: func() any { return []byte(gofakeit.SentenceSimple()) },
				get:         func(et *EMVProprietaryTemplate) any { return et.CDOL1 },
				set:         func(et *EMVProprietaryTemplate, val any) { et.CDOL1 = val.(tlv.TL) },
			},
			{
				field:       "CDOL2",
				emptyValue:  func() any { return "" },
				randomValue: func() any { return gofakeit.SentenceSimple() },
				get:         func(et *EMVProprietaryTemplate) any { return et.CDOL2 },
				set:         func(et *EMVProprietaryTemplate, val any) { et.CDOL2 = val.(tlv.TL) },
			},
			{
				field:       "VersionNumber1",
				emptyValue:  func() any { return "" },
				randomValue: func() any { return gofakeit.SentenceSimple() },
				get:         func(et *EMVProprietaryTemplate) any { return et.VersionNumber1 },
				set:         func(et *EMVProprietaryTemplate, val any) { et.VersionNumber1 = val.(string) },
			},
			{
				field:       "ICCPublicKeyCertificate",
				emptyValue:  func() any { return "" },
				randomValue: func() any { return gofakeit.SentenceSimple() },
				get:         func(et *EMVProprietaryTemplate) any { return et.ICCPublicKeyCertificate },
				set:         func(et *EMVProprietaryTemplate, val any) { et.ICCPublicKeyCertificate = val.(string) },
			},
			{
				field:       "ICCPublicKeyExponent",
				emptyValue:  func() any { return "" },
				randomValue: func() any { return gofakeit.SentenceSimple() },
				get:         func(et *EMVProprietaryTemplate) any { return et.ICCPublicKeyExponent },
				set:         func(et *EMVProprietaryTemplate, val any) { et.ICCPublicKeyExponent = val.(string) },
			},
			{
				field:       "DDOL",
				emptyValue:  func() any { return "" },
				randomValue: func() any { return gofakeit.SentenceSimple() },
				get:         func(et *EMVProprietaryTemplate) any { return et.DDOL },
				set:         func(et *EMVProprietaryTemplate, val any) { et.DDOL = val.(string) },
			},
			{
				field:       "CVMListBytes",
				emptyValue:  func() any { return "" },
				randomValue: func() any { return gofakeit.SentenceSimple() },
				get:         func(et *EMVProprietaryTemplate) any { return et.CVMListBytes },
				set:         func(et *EMVProprietaryTemplate, val any) { et.CVMListBytes = val.([]byte) },
			},
		}
		for _, tc := range testCases {
			t.Run(tc.field, func(t *testing.T) {
				t.Run("When value is set to field of A", func(t *testing.T) {
					var a EMVProprietaryTemplate
					gofakeit.Struct(&a)
					var b EMVProprietaryTemplate
					gofakeit.Struct(&b)

					fakeValue := tc.randomValue()
					tc.set(&a, fakeValue)

					merged := a.Merge(b)
					mergedFieldValue := tc.get(&merged)
					assert.Equal(t, fakeValue, mergedFieldValue)
				})
			})
			t.Run(tc.field, func(t *testing.T) {
				t.Run("When value of A is empty and value is set to field of B", func(t *testing.T) {
					var a EMVProprietaryTemplate
					gofakeit.Struct(&a)
					var b EMVProprietaryTemplate
					gofakeit.Struct(&b)

					tc.set(&a, tc.emptyValue())
					fakeValue := tc.randomValue()
					tc.set(&b, fakeValue)

					merged := a.Merge(b)
					mergedFieldValue := tc.get(&merged)
					assert.Equal(t, fakeValue, mergedFieldValue)
				})
			})
		}
	})
}
