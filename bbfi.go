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
			addPtr(buf, 1)
		case MINUS:
			subPtr(buf, 1)
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
		}
		inst = inst.next
	}
	return nil
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
	execute(inst, buf)
}
