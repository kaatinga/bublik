package bublyk

import "testing"

func BenchmarkWithAssets(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		maximumDate.Format(PostgreSQLFormat)
	}
}

func BenchmarkUsingTime(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		maximumDate.Format("2006_01-02")
	}
}
