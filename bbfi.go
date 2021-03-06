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
	MOV
	MOV_NONZERO
	NOP
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

func movePtrByOffset(buf *Buffer, offset int) *Buffer {
	if offset >= 0 {
		if buf.bidx+offset > BUF_MAX-1 {
			diff := BUF_MAX - 1 - buf.bidx
			buf.bidx = BUF_MAX - 1
			if buf.next == nil {
				new_buf := appendNewBuffer(buf)
				new_buf.bidx = offset - diff
				buf = new_buf
			} else {
				buf = buf.next
				buf.bidx = offset - diff
			}
			return buf
		} else {
			buf.bidx += offset
			return buf
		}
	} else {
		if buf.bidx-offset < 0 {
			if buf.prev == nil {
				return nil
			}
			diff := offset + buf.bidx
			buf = buf.prev
			buf.bidx += diff
			return buf
		} else {
			buf.bidx += offset
			return buf
		}
	}
}

func movePtrToNonzero(buf *Buffer, gap int) *Buffer {
	for i := buf.bidx; ; i += gap {
		if i > BUF_MAX-1 {
			buf.bidx = BUF_MAX - 1
			diff := i - BUF_MAX - 1
			if buf.next == nil {
				new_buf := appendNewBuffer(buf)
				new_buf.bidx = diff
				buf = new_buf
			} else {
				buf = buf.next
				buf.bidx = diff
			}
			i = diff
		} else if i < 0 {
			buf.bidx = 0
			buf = buf.prev
			buf.bidx += i
			i = buf.bidx
		}

		if buf.buf[i] == 0 {
			buf.bidx = i
			return buf
		}
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
		case MOV:
			buf = movePtrByOffset(buf, inst.inst_parameter)
		case MOV_NONZERO:
			buf = movePtrToNonzero(buf, inst.inst_parameter)
		case NOP:
			break
		}
		inst = inst.next
	}
	return nil
}

func optimizeConsecutiveArithmetic(inst *Instruction) {
	// Transform consecutive arithmetic operation, ex: ++++ to PLUS 4
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
			if start != nil && start != inst.prev {
				if sum > 0 {
					start.inst_type = PLUS
				} else if sum < 0 {
					start.inst_type = MINUS
					sum = -sum
				}
				start.inst_parameter = sum

				sum = 0
				start.next = inst
				inst.prev = start
			}
			start = nil
		}
		inst = inst.next
	}

	if start != nil {
		start.next = nil
	}
}

func optimizeConsecutiveMovement(inst *Instruction) {
	// Transform consecutive pointer movement, ex: <<<< to MOV -4
	offset := 0
	start := (*Instruction)(nil)
	for inst != nil {
		if inst.inst_type == PREV || inst.inst_type == NEXT {
			if start == nil {
				if inst.inst_type == NEXT {
					offset = 1
				} else {
					offset = -1
				}
				start = inst
			} else {
				if inst.inst_type == NEXT {
					offset += 1
				} else {
					offset -= 1
				}
			}
		} else {
			if start != nil && start != inst.prev {
				start.inst_type = MOV
				start.inst_parameter = offset

				start.next = inst
				inst.prev = start
			}
			offset = 0
			start = nil
		}
		inst = inst.next
	}

	if start != nil {
		start.next = nil
	}
}

func optimizeSubstractionToZero(inst *Instruction) {
	// For pattern: [-]
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
					if new_inst.next != nil {
						new_inst.next.prev = new_inst
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

func optimizeMovementLoop(inst *Instruction) {
	// For pattern: [>], [<], [>>>]
	for inst != nil {
		if inst.inst_type == LOOP {
			inner_inst := inst.next
			if inner_inst.next.inst_type == END {
				is_movement := false
				switch inner_inst.inst_type {
				case NEXT:
					is_movement = true
					inst.inst_type = MOV_NONZERO
					inst.inst_parameter = 1
				case PREV:
					is_movement = true
					inst.inst_type = MOV_NONZERO
					inst.inst_parameter = -1
				case MOV:
					is_movement = true
					inst.inst_type = MOV_NONZERO
					inst.inst_parameter = inner_inst.inst_parameter
				}
				if is_movement {
					inst.next = inst.inst_pair.next
					if inst.next != nil {
						inst.next.prev = inst
					}
					inst.inst_pair = nil
				}
			}
		}
		inst = inst.next
	}
}

func optimizeAssignment(inst *Instruction) {
	// For pattern: [-]++++++, considered as variable assignment
	for inst != nil {
		if inst.inst_type == SET {
			if inst.next != nil {
				if inst.next.inst_type == PLUS {
					inst.inst_parameter = inst.next.inst_parameter
					inst.next = inst.next.next
					if inst.next != nil {
						inst.next.prev = inst
					}
				} else if inst.inst_type == MINUS {
					inst.inst_parameter = 0 - inst.next.inst_parameter
					inst.next = inst.next.next
				}
			}
		}
		inst = inst.next
	}
}

func optimizeSimpleLoopMultiplication(inst *Instruction) {
	// For pattern: ++++++[->+>++++<<]
	for inst != nil {
		if inst.inst_type == LOOP {
			// Determine if this loop represents the pattern
			if inst.inst_parameter == 1 || inst.prev == nil {
				// Has subloops
				inst = inst.next
				continue
			}
			if inst.prev != nil && inst.prev.inst_type != PLUS && inst.prev.inst_type != SET {
				inst = inst.next
				continue
			}
			loop_inst := inst.next
			pattern_match := true
			offset := 0
			for loop_inst != inst.inst_pair {
				switch loop_inst.inst_type {
				case NEXT:
					offset += 1
				case PREV:
					offset -= 1
				case MOV:
					offset += loop_inst.inst_parameter
				case READ:
					pattern_match = false
				}
				loop_inst = loop_inst.next
			}
			if offset != 0 {
				pattern_match = false
			}
			if pattern_match {
				mul_factor := inst.prev.inst_parameter
				loop_inst = inst.next
				for loop_inst != inst.inst_pair {
					if loop_inst.inst_type == PLUS {
						loop_inst.inst_parameter *= mul_factor
					} else if loop_inst.inst_type == MINUS && loop_inst.inst_parameter == 1 {
						// Check if it's substraction to condition index
						if loop_inst.prev.inst_type == LOOP || loop_inst.next.inst_type == END {
							loop_inst.inst_type = NOP
						}
					} else if loop_inst.inst_type == MINUS {
						loop_inst.inst_parameter *= mul_factor
					}
					loop_inst = loop_inst.next
				}
				inst.prev.inst_type = NOP
				inst.prev.inst_parameter = 0
				inst.prev.next = inst.next
				inst.next.prev = inst.prev
				loop_inst.prev.next = loop_inst.next
				inst = loop_inst
			} else {
				inst = inst.inst_pair
			}
		}
		inst = inst.next
	}
}

func optimizeInstructions(inst *Instruction) *Instruction {
	head := inst
	optimizeConsecutiveArithmetic(inst)
	optimizeConsecutiveMovement(inst)
	optimizeSubstractionToZero(inst)
	optimizeMovementLoop(inst)
	optimizeAssignment(inst)
	optimizeSimpleLoopMultiplication(inst)
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
			for e := stack.Front(); e != nil; e = e.Next() {
				e.Value.(*Instruction).inst_parameter = 1
			}
		}
		i, _ = fmt.Scanf("%c", &c)
	}
	inst = optimizeInstructions(inst)
	execute(inst, buf)
}
