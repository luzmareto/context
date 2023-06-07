package course_context

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	/* Context.Background di awal, biasanya kosong
	dibawah ini adalah basic penggunaan Context
	*/
	background := context.Background()
	fmt.Println(background)

	todo := context.TODO()
	fmt.Println(todo)
}

/*
DIAGRAM CONTEXT adalah hirarki
Context.Background di awal, biasanya kosong
Context pertama bisa melahirkan Context kedua dengan value dari Context pertama
Context kedua bisa melahirkan Context ketiga dengan value yang sama dari context kedua
Namun Context pertama tidak mendapatkan value dari Context ketiga
*/

func TestContextWithValue(t *testing.T) {
	contextA := context.Background()

	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")

	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	contextF := context.WithValue(contextC, "f", "F")
	contextG := context.WithValue(contextF, "g", "G")

	//maka tiap context memiliki parent yaitu context.Background
	fmt.Println(contextA) //context.Background
	fmt.Println(contextB) //context.Background.WithValue(type string, val B)
	fmt.Println(contextC) //context.Background.WithValue(type string, val C)
	fmt.Println(contextD) //context.Background.WithValue(type string, val B).WithValue(type string, val E)
	fmt.Println(contextE) //context.Background.WithValue(type string, val B).WithValue(type string, val E)
	fmt.Println(contextF) //context.Background.WithValue(type string, val C).WithValue(type string, val F)
	fmt.Println(contextG) //context.Background.WithValue(type string, val C).WithValue(type string, val F).WithValue(type string, val G)

	//menanyakan value ke parents
	fmt.Println(contextF.Value("f")) //F
	fmt.Println(contextF.Value("c")) //C
	fmt.Println(contextF.Value("b")) //nil
	fmt.Println(contextA.Value("b")) //nil
}

func CreateCounter(ctx context.Context) chan int {
	//var penampung
	destination := make(chan int)

	//anonym func
	go func() {
		defer close(destination) //menghentikan destination
		counter := 1
		for {
			select {
			case <-ctx.Done():
				return //menghentikan Anonym Func
			default:
				destination <- counter
				counter++
				time.Sleep(1 * time.Second)
			}
		}
	}()

	return destination //return dari chan
}

func TestContextWithCancel(t *testing.T) {
	//mengecek jumlah Goroutine
	fmt.Println("total Goroutine", runtime.NumGoroutine()) //2

	//membuat context parent
	parent := context.Background()

	//
	ctx, cancel := context.WithCancel(parent)

	//memanggil func Createcounter
	destination := CreateCounter(ctx)

	//perulangan 10x
	for n := range destination {
		fmt.Println("counter", n)
		if n == 10 { //mencetak sebanyak 10x
			break
		}
	}

	//mengirim sinyal cancel ke context
	cancel()

	time.Sleep(2 * time.Second)

	fmt.Println("total Goroutine", runtime.NumGoroutine()) //3
}

func TestContextWithTimeOut(t *testing.T) {
	//mengecek jumlah Goroutine
	fmt.Println("total Goroutine", runtime.NumGoroutine()) //2

	//membuat context parent
	parent := context.Background()

	//proses berhenti pada detike ke 5, meskipun masih ada tugasnya
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()

	//memanggil func Createcounter
	destination := CreateCounter(ctx)
	fmt.Println("total Goroutine", runtime.NumGoroutine())

	//perulangan tanpa henti
	for n := range destination {
		fmt.Println("counter", n)
	}

	time.Sleep(2 * time.Second)

	fmt.Println("total Goroutine", runtime.NumGoroutine()) //3
}

func TestContextWithDeadLine(t *testing.T) {
	//mengecek jumlah Goroutine
	fmt.Println("total Goroutine", runtime.NumGoroutine()) //2

	//membuat context parent
	parent := context.Background()

	//proses berhenti pada detike ke 5, meskipun masih ada tugasnya
	ctx, cancel := context.WithDeadline(parent, time.Now().Add(5*time.Second))
	defer cancel()

	//memanggil func Createcounter
	destination := CreateCounter(ctx)
	fmt.Println("total Goroutine", runtime.NumGoroutine())

	//perulangan tanpa henti
	for n := range destination {
		fmt.Println("counter", n)
	}

	time.Sleep(2 * time.Second)

	fmt.Println("total Goroutine", runtime.NumGoroutine()) //3
}
