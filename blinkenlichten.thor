require 'serialport'
require 'pry'

class Blinkenlichten < Thor
  class_option :device, type: :string, default: Dir["/dev/tty.usbmodem*"].first
  class_option :baud, type: :string, default: 9600

  desc "on", "turn on all the lights"
  def on
    port = init_port(options)
    leds = num_leds(port)
    (0...leds).each do |i|
      port.write("s#{i} 255 255 255\n")
    end
  end

  desc "off", "turn off all the lights"
  def off
    port = init_port(options)
    leds = num_leds(port)
    (0...leds).each do |i|
      port.write("s#{i} 0 0 0\n")
    end
  end

  desc "console", "open a pry console with the port open"
  def console
    port = init_port(options)
    binding.pry
  end

  private
  def init_port(options)
    port = SerialPort.new(options[:device], baud: options[:baud], dtr: 0)
    loop do
      break if port.gets.chomp == "ready."
    end
    port
  end

  def num_leds(port)
    port.write("c\n")
    port.gets.chomp.to_i
  end
end
