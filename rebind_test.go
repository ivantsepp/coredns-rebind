package rebind

import (
	"net"
	"testing"
)

func TestGetIPAnswerFirstThenSecond(t *testing.T) {
	firstIP := net.ParseIP("1.2.3.4")
	secondIP := net.ParseIP("0.0.0.0")

	r := Rebind{
		firstIP:  firstIP,
		secondIP: secondIP,
		strategy: firstThenSecondStrategy,
		visited:  make(map[string]int),
	}

	if ans := r.getIPAnswerFirstThenSecond("example.com"); !ans.Equal(firstIP) {
		t.Errorf("Expected %v, got %v", firstIP, ans)
	}

	if ans := r.getIPAnswerFirstThenSecond("example.com"); !ans.Equal(secondIP) {
		t.Errorf("Expected %v, got %v", secondIP, ans)
	}

	if ans := r.getIPAnswerFirstThenSecond("example.com"); !ans.Equal(secondIP) {
		t.Errorf("Expected %v, got %v", secondIP, ans)
	}
}

func TestGetIPAnswerRoundRobin(t *testing.T) {
	firstIP := net.ParseIP("1.2.3.4")
	secondIP := net.ParseIP("0.0.0.0")

	r := Rebind{
		firstIP:  firstIP,
		secondIP: secondIP,
		strategy: roundRobinStrategy,
		visited:  make(map[string]int),
	}

	if ans := r.getIPAnswerRoundRobin("example.com"); !ans.Equal(firstIP) {
		t.Errorf("Expected %v, got %v", firstIP, ans)
	}

	if ans := r.getIPAnswerRoundRobin("example.com"); !ans.Equal(secondIP) {
		t.Errorf("Expected %v, got %v", secondIP, ans)
	}

	if ans := r.getIPAnswerRoundRobin("example.com"); !ans.Equal(firstIP) {
		t.Errorf("Expected %v, got %v", firstIP, ans)
	}
}

func TestGetIPAnswerRandom(t *testing.T) {
	firstIP := net.ParseIP("1.2.3.4")
	secondIP := net.ParseIP("0.0.0.0")

	r := Rebind{
		firstIP:  firstIP,
		secondIP: secondIP,
		strategy: randomStrategy,
		visited:  make(map[string]int),
	}

	if ans := r.getIPAnswerRandom(); !ans.Equal(firstIP) && !ans.Equal(secondIP) {
		t.Errorf("Expected %v or %v, got %v", firstIP, secondIP, ans)
	}

	if ans := r.getIPAnswerRandom(); !ans.Equal(firstIP) && !ans.Equal(secondIP) {
		t.Errorf("Expected %v or %v, got %v", firstIP, secondIP, ans)
	}

	if ans := r.getIPAnswerRandom(); !ans.Equal(firstIP) && !ans.Equal(secondIP) {
		t.Errorf("Expected %v or %v, got %v", firstIP, secondIP, ans)
	}
}
