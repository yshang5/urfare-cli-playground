package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "mytool",
		Usage: "Basic CLI demo",
		Commands: []*cli.Command{
			{
				Name:  "hello",
				Usage: "Say hello",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("Hello!")
					return nil
				},
			},

			{
				Name:  "calculator",
				Usage: "get result for a math expression, valid operators including + - * / %",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "expression",
						Usage: "math expression",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					exp := cmd.String("expression")
					i, err := calculate(exp)
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Println(i)
					}
					return nil
				},
			},
		},
	}

	app.Run(context.Background(), os.Args)
}

var priorityMap = map[string]int{
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
	"(": 0,
	")": 0,
	"%": 2,
}
var err = errors.New("无效的表达式")

func calculate(expression string) (int, error) {
	numStack := Stack[int]{}
	optStack := Stack[string]{}

	numStr := ""
	for _, ch := range expression {
		s := string(ch)
		if ch >= '0' && ch <= '9' {
			numStr += s
		} else if _, ok := priorityMap[s]; ok {
			switch s {
			case "(":
				// 左括号不需要转换numStr，直接压入运算符栈
				optStack.Push(s)
			case "+", "-", "*", "/", "%":
				// 遇到运算符，先把之前累积的数字压栈
				if len(numStr) > 0 {
					num, _ := strconv.Atoi(numStr)
					numStack.Push(num)
					numStr = ""
				}
				err := dealOperators(&numStack, &optStack, s)
				if err != nil {
					return 0, err
				}
			case ")":
				// 右括号前可能有数字，先处理数字
				if len(numStr) > 0 {
					num, _ := strconv.Atoi(numStr)
					numStack.Push(num)
					numStr = ""
				}
				err = dealParentheses(&numStack, &optStack)
				if err != nil {
					return 0, err
				}
			}
		} else {
			return 0, err
		}
	}

	// 处理最后剩余的数字（如果有）
	if len(numStr) > 0 {
		num, _ := strconv.Atoi(numStr)
		numStack.Push(num)
	}

	for optStack.Len() != 0 {
		b, bExists := numStack.Pop()
		if !bExists {
			return 0, err
		}
		a, aExists := numStack.Pop()
		if !aExists {
			return 0, err
		}
		opt, _ := optStack.Pop()
		res := cal(a, b, opt)
		numStack.Push(res)
	}
	if numStack.Len() == 1 {
		ans, _ := numStack.Pop()
		return ans, nil
	}
	return 0, err
}

func dealParentheses(nums *Stack[int], opts *Stack[string]) error {
	opt, exists := opts.Pop()
	for exists && opt != "(" {
		b, bExists := nums.Pop()
		if !bExists {
			return errors.New("无效的表达式")
		}
		a, aExists := nums.Pop()
		if !aExists {
			return errors.New("无效的表达式")
		}
		res := cal(a, b, opt)
		nums.Push(res)
		opt, exists = opts.Pop()
	}
	if !exists {
		return err
	}

	return nil
}

func dealOperators(nums *Stack[int], opts *Stack[string], opt string) error {
	lastOpt, exist := opts.Peek()
	curPriority := priorityMap[opt]
	lastPriority := priorityMap[lastOpt]
	//上一个运算符大于当前运算符，就先算
	if exist && lastPriority >= curPriority {
		opts.Pop() // 先弹出要计算的运算符
		b, bExists := nums.Pop()
		if !bExists {
			return errors.New("无效的表达式")
		}
		a, aExists := nums.Pop()
		if !aExists {
			return errors.New("无效的表达式")
		}
		res := cal(a, b, lastOpt)
		nums.Push(res)
	}
	opts.Push(opt)
	return nil
}

func cal(a int, b int, opt string) int {
	var res int
	switch opt {
	case "+":
		res = a + b
		break
	case "-":
		res = a - b
		break

	case "*":
		res = a * b
		break
	case "/":
		res = a / b
		break
	case "%":
		res = a % b
		break
	}
	return res
}

type Stack[T any] struct {
	data []T
}

func (s *Stack[T]) Push(v T) {
	s.data = append(s.data, v)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	v := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return v, true
}

func (s *Stack[T]) Peek() (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	return s.data[len(s.data)-1], true
}
func (s *Stack[T]) Len() int {
	return len(s.data)
}
