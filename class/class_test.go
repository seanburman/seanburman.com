package class

import "testing"

func TestMain(m *testing.M) {
	m.Run()
}

func TestClass_String(t *testing.T) {
	tests := []struct {
		name string
		c    Class
		want string
	}{
		{
			name: "Test Class String",
			c:    "bg-white p-8 rounded-lg shadow-lg w-96",
			want: "bg-white p-8 rounded-lg shadow-lg w-96",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("Class.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClass_Add(t *testing.T) {
	tests := []struct {
		name  string
		c     Class
		class []Class
		want  Class
	}{
		{
			name:  "Test Class Add",
			c:     "bg-white p-8 rounded-lg shadow-lg w-96",
			class: []Class{"bg-black", "hover:bg-gray-700"},
			want:  "bg-white p-8 rounded-lg shadow-lg w-96 bg-black hover:bg-gray-700",
		},
		{
			name:  "Test Class Add",
			c:     "bg-white p-8 rounded-lg shadow-lg w-96",
			class: []Class{"bg-black", "hover:bg-gray-700", "bg-white", "p-8", "rounded-lg"},
			want:  "bg-white p-8 rounded-lg shadow-lg w-96 bg-black hover:bg-gray-700",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Add(tt.class...); got != tt.want {
				t.Errorf("Class.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClass_Remove(t *testing.T) {
	tests := []struct {
		name  string
		c     Class
		class Class
		want  Class
	}{
		{
			name:  "Test Class Remove",
			c:     "bg-white p-8 rounded-lg shadow-lg w-96 bg-black hover:bg-gray-700",
			class: "bg-black",
			want:  "bg-white p-8 rounded-lg shadow-lg w-96  hover:bg-gray-700",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Remove(tt.class); got != tt.want {
				t.Errorf("Class.Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}
