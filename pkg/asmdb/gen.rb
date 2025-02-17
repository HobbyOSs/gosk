require 'rexml/document'
require 'yaml'


doc = REXML::Document.new(File.open("x86reference/x86reference.xml"))
data = []

def build_operand_string_exp(addr, type)
  addr_s = ""

  case addr
  when 'A'
    addr_s = "ptr"
  when 'BA', 'BB', 'BD', 'M', 'X', 'Y'
    addr_s = "m"
  when 'C'
    addr_s = "CRn"
  when 'D'
    addr_s = "DRn"
  when 'E'
    addr_s = "r/m"
  when 'ES'
    addr_s = "STi/m"
  when 'EST'
    addr_s = "STi"
  when 'F', 'SC'
    addr_s = "-"
  when 'G', 'H', 'R', 'Z'
    addr_s = "r"
  when 'I'
    addr_s = "imm"
  when 'J'
    addr_s = "rel"
  when 'N', 'P'
    addr_s = "mm"
  when '0'
    addr_s = "moffs"
  when 'Q'
    addr_s = "mm/m64"
  when 'S'
    addr_s = "Sreg"
  when 'T'
    addr_s = "TRn"
  when 'U', 'V'
    addr_s = "xmm"
  when 'W'
    addr_s = "xmm/m"
  end

  case type
  when 'a'
    addr_s += "16/32&16/32"
  when 'b', 'bs', 'bss'
    addr_s += "8"
  when 'bcd'
    addr_s += "80dec"
  when 'bsq', 'pt', 's', 'sd', 'ss', 't'
    addr_s += "-"
  when 'c', 'si'
    addr_s += "?"
  when 'd'
    addr_s += "32"
  when 'di'
    addr_s += "32int"
  when 'dq'
    addr_s += "128"
  when 'dqp'
    addr_s += "32/64"
  when 'dr'
    addr_s += "64real"
  when 'ds'
    addr_s += "32"
  when 'e'
    addr_s += "14/28"
  when 'er'
    addr_s += "80real"
  when 'p'
    addr_s += "16:16/32"
  when 'pi'
    addr_s += "(64)"
  when 'pd'
    addr_s += ""
  when 'ps'
    addr_s += "(128)"
  when 'psq'
    addr_s += "64"
  when 'ptp'
    addr_s += "16:16/32/64"
  when 'q', 'qp'
    addr_s += "64"
  when 'qi'
    addr_s += "64int"
  when 'sr'
    addr_s += "32real"
  when 'st'
    addr_s += "94/108"
  when 'stx'
    addr_s += "512"
  when 'v', 'vds', 'vs'
    addr_s += "16/32"
  when 'vq'
    addr_s += "64/16"
  when 'vqp'
    addr_s += "16/32/64"
  when 'w'
    addr_s += "16"
  when 'wi'
    addr_s += "16int"
  end

  return addr_s
end



REXML::XPath.match(doc, "/x86reference/one-byte/pri_opcd/entry | /x86reference/two-byte/pri_opcd/entry").each do |entry|
  if entry.elements["syntax/mnem"]
    opcd = entry.parent.attributes["value"]
    mnem = ""
    operands = []

    entry.elements.each('syntax/*') do |elem|
      case elem.name
      when 'mnem'
        mnem = elem&.text
      when 'src', 'dst'
        # If an operand is set up using italic, it is an implicit operand,
        # which is not explicitly used. If an operand is set up using boldface,
        # it is modified by the instruction.
        # つまり、上記の場合実際のアセンブラの引数にはならない
        if elem.attributes["displayed"] == 'no'
          operands << nil
        elsif elem.text.nil?
          addr = elem.elements['a']&.text
          type = elem.elements['t']&.text

          operands << {
            "#{elem.name}" => {
              "operand_s" => build_operand_string_exp(addr, type),
              "a" => addr,
              "t" => type
            }
          }
        else
          operands << {
            "#{elem.name}" => {
              "operand_s" => elem.text
            }
          }
        end
      end
    end

    op1 = operands[0]
    op2 = operands[1]

    # 00: 8086
    # 01: 80186
    # 02: 80286
    # 03: 80386
    # 04: 80486
    # P1 (05): Pentium (1)
    # PX (06): Pentium with MMX
    # PP (07): Pentium Pro
    # P2 (08): Pentium II
    # P3 (09): Pentium III
    # P4 (10): Pentium 4
    # C1 (11): Core (1)
    # C2 (12): Core 2
    # C7 (13): Core i7
    # IT (99): Itanium (only geek editions)
    proc_s = entry.elements["proc_start"]&.text
    proc_s = proc_s || "00+"
    # 説明
    desc = entry.elements["note/brief"].text
    data << {
      "mnem" => "#{mnem}",
      "opcd" => "#{opcd}" ,
      "op1" => op1,
      "op2" => op2,
      "proc" => proc_s,
      "desc" => "#{desc}"
    }
  end
end

YAML.dump(data, File.open('x86reference.yml', 'w'))
