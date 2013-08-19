package main

import (
	"container/list"
	"fmt"
)

const (
	PREV = iota + 13
	NEXT
	PLUS
	MINUS
	PRINT
	READ
	LOOP
	END
	SET
)

const (
	BUF_MAX = iota + 20000
)

type Buffer struct {
	buf  [BUF_MAX]int
	bidx int
	next *Buffer
	prev *Buffer
}

type Instruction struct {
	inst_type      int
	inst_parameter int
	inst_pair      *Instruction
	next           *Instruction
	prev           *Instruction
}

func appendNewBuffer(head *Buffer) *Buffer {
	for head.next != nil {
		head = head.next
	}
	head.next = new(Buffer)
	head.next.prev = head
	return head.next
}

func appendNewInstruction(head *Instruction, inst_type int, inst_parameter int, inst_pair *Instruction) *Instruction {
	if head.inst_type == 0 { // The very first instruction
		head.inst_type = inst_type
		head.inst_parameter = inst_parameter
		return head
	} else {
		for head.next != nil {
			head = head.next
		}
		head.next = new(Instruction)
		head.next.prev = head
		head.next.inst_type = inst_type
		head.next.inst_parameter = inst_parameter
		head.next.inst_pair = inst_pair
		if inst_pair != nil {
			inst_pair.inst_pair = head.next // Since `]` will always appear in the end of instruction list
		}
		return head.next
	}
}

func movePtrToPrev(buf *Buffer) *Buffer {
	if buf.bidx == 0 {
		if buf.prev == nil {
			return nil
		}
		return buf.prev
	} else {
		buf.bidx -= 1
		return buf
	}
}

func movePtrToNext(buf *Buffer) *Buffer {
	if buf.bidx == BUF_MAX-1 {
		if buf.next == nil {
			return appendNewBuffer(buf)
		} else {
			return buf.next
		}
	} else {
		buf.bidx += 1
		return buf
	}
}

func addPtr(buf *Buffer, val int) {
	buf.buf[buf.bidx] += val
}

func subPtr(buf *Buffer, val int) {
	buf.buf[buf.bidx] -= val
}

func printPtr(buf *Buffer) {
	fmt.Printf("%c", buf.buf[buf.bidx])
}

func readPtr(buf *Buffer) {
	fmt.Scanf("%c", &(buf.buf[buf.bidx]))
}

func execute(inst *Instruction, buf *Buffer) *Instruction {
	for inst != nil {
		switch inst.inst_type {
		default:
			return nil
		case PREV:
			buf = movePtrToPrev(buf)
		case NEXT:
			buf = movePtrToNext(buf)
		case PLUS:
			addPtr(buf, inst.inst_parameter)
		case MINUS:
			subPtr(buf, inst.inst_parameter)
		case PRINT:
			printPtr(buf)
		case READ:
			readPtr(buf)
		case LOOP:
			next_inst := inst.inst_pair
			for buf.buf[buf.bidx] != 0 {
				execute(inst.next, buf)
			}
			inst = next_inst
		case END:
			return inst
		case SET:
			buf.buf[buf.bidx] = inst.inst_parameter
		}
		inst = inst.next
	}
	return nil
}

func optimizeConsecutiveArithmetic(inst *Instruction) {
	sum := 0
	start := (*Instruction)(nil)
	for inst != nil {
		if inst.inst_type == PLUS || inst.inst_type == MINUS {
			if start == nil {
				if inst.inst_type == PLUS {
					sum = 1
				} else {
					sum = -1
				}
				start = inst
			} else {
				if inst.inst_type == PLUS {
					sum += 1
				} else {
					sum -= 1
				}
			}
		} else {
			if start != nil {
				if sum > 0 {
					start.inst_type = PLUS
				} else if sum < 0 {
					start.inst_type = MINUS
					sum = -sum
				}
				start.inst_parameter = sum

				sum = 0
				start.next = inst
			}
			start = nil
		}
		inst = inst.next
	}

	if start != nil {
		start.next = nil
	}
}

func optimizeSubstractionToZero(inst *Instruction) {
	for inst != nil {
		if inst.inst_type == LOOP {
			if inst.next != nil && inst.next.next != nil {
				if inst.next.inst_type == MINUS && inst.next.next.inst_type == END {
					new_inst := new(Instruction)
					new_inst.prev = inst.prev
					if inst.prev != nil {
						inst.prev.next = new_inst
					}
					new_inst.next = inst.next.next.next
					if inst.next.next.next != nil {
						inst.next.next.next.prev = new_inst
					}
					new_inst.inst_type = SET
					new_inst.inst_parameter = 0
					inst = new_inst
				}
			}
		}
		inst = inst.next
	}
}

func optimizeInstructions(inst *Instruction) *Instruction {
	head := inst
	optimizeConsecutiveArithmetic(inst)
	optimizeSubstractionToZero(inst)
	return head
}

func main() {
	inst := new(Instruction)
	buf := new(Buffer)
	stack := list.New()
	var c int
	i, _ := fmt.Scanf("%c", &c)
	for i != 0 {
		switch c {
		case '>':
			appendNewInstruction(inst, NEXT, 0, nil)
		case '<':
			appendNewInstruction(inst, PREV, 0, nil)
		case '+':
			appendNewInstruction(inst, PLUS, 1, nil)
		case '-':
			appendNewInstruction(inst, MINUS, 1, nil)
		case '.':
			appendNewInstruction(inst, PRINT, 0, nil)
		case ',':
			appendNewInstruction(inst, READ, 0, nil)
		case '[':
			begin := appendNewInstruction(inst, LOOP, 0, nil)
			stack.PushBack(begin)
		case ']':
			begin := stack.Back().Value.(*Instruction)
			appendNewInstruction(inst, END, 0, begin)
			stack.Remove(stack.Back())
		}
		i, _ = fmt.Scanf("%c", &c)
	}
	inst = optimizeInstructions(inst)
	execute(inst, buf)
}
