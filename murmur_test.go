package murmur3

import "testing"

func TestMurmurHash3_Sum128(t *testing.T) {
	type fields struct {
		k1 uint64
		k2 uint64
		h1 uint64
		h2 uint64
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint64
		want1  uint64
	}{
		{
			name: "123456",
			fields: fields{
				h1: 123456,
				h2: 123456,
			},
			args: args{
				data: []byte("helloqwertyuiophelloqwertyuiop"),
			},
			want:  14288351280354662,
			want1: 2696564043591294180,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MurmurHash3{
				k1: tt.fields.k1,
				k2: tt.fields.k2,
				h1: tt.fields.h1,
				h2: tt.fields.h2,
			}
			got, got1 := m.Sum128(tt.args.data)
			if got != tt.want {
				t.Errorf("MurmurHash3.Sum128() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("MurmurHash3.Sum128() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMurmurHash3_k1Calc(t *testing.T) {
	type fields struct {
		k1 uint64
		k2 uint64
		h1 uint64
		h2 uint64
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "123456",
			fields: fields{
				h1: 123456,
				h2: 123456,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MurmurHash3{
				k1: tt.fields.k1,
				k2: tt.fields.k2,
				h1: tt.fields.h1,
				h2: tt.fields.h2,
			}
			m.k1Calc()
		})
	}
}

func TestMurmurHash3_k2Calc(t *testing.T) {
	type fields struct {
		k1 uint64
		k2 uint64
		h1 uint64
		h2 uint64
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "123456",
			fields: fields{
				h1: 123456,
				h2: 123456,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MurmurHash3{
				k1: tt.fields.k1,
				k2: tt.fields.k2,
				h1: tt.fields.h1,
				h2: tt.fields.h2,
			}
			m.k2Calc()
		})
	}
}

func Test_rotl64(t *testing.T) {
	type args struct {
		x uint64
		r uint8
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "10 and 8",
			args: args{
				x: 10,
				r: 8,
			},
			want: 2560,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rotl64(tt.args.x, tt.args.r); got != tt.want {
				t.Errorf("rotl64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fmix64(t *testing.T) {
	type args struct {
		x uint64
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "10",
			args: args{
				x: 10,
			},
			want: 7233188113542599437,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fmix64(tt.args.x); got != tt.want {
				t.Errorf("fmix64() = %v, want %v", got, tt.want)
			}
		})
	}
}
