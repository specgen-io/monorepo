=begin

source: https://github.com/jackc/tod

Copyright (c) 2010-2015 Jack Christensen

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

=end

module [[.ModuleName]]
  class TimeOfDay
    include Comparable

    def self.jsoner_deserialize(json_value)
      TimeOfDay.parse(T.check(String, json_value))
    end
    def self.jsoner_serialize(value)
      T.check(TimeOfDay, value).to_s
    end

    attr_reader :hour, :minute, :second, :second_of_day
    alias_method :min, :minute
    alias_method :sec, :second
    alias_method :to_i, :second_of_day

    PARSE_24H_REGEX = /
      \A
      ([01]?\d|2[0-4])
      :?
      ([0-5]\d)?
      :?
      ([0-5]\d)?
      \z
    /x

    PARSE_12H_REGEX = /
      \A
      (0?\d|1[0-2])
      :?
      ([0-5]\d)?
      :?
      ([0-5]\d)?
      \s*
      ([ap])
      \.?
      \s*
      m?
      \.?
      \z
    /x

    WORDS = {
        "noon" => "12pm".freeze,
        "midnight" => "12am".freeze
    }

    NUM_SECONDS_IN_DAY = 86400
    NUM_SECONDS_IN_HOUR = 3600
    NUM_SECONDS_IN_MINUTE = 60

    FORMATS = {
        short: "%-l:%M %P".freeze,
        medium: "%-l:%M:%S %P".freeze,
        time: "%H:%M".freeze
    }

    def initialize(h, m=0, s=0)
      @hour = Integer(h)
      @minute = Integer(m)
      @second = Integer(s)

      raise ArgumentError, "hour must be between 0 and 24" unless (0..24).include?(@hour)
      if @hour == 24 && (@minute != 0 || @second != 0)
        raise ArgumentError, "hour can only be 24 when minute and second are 0"
      end
      raise ArgumentError, "minute must be between 0 and 59" unless (0..59).include?(@minute)
      raise ArgumentError, "second must be between 0 and 59" unless (0..59).include?(@second)

      @second_of_day = @hour * 60 * 60 + @minute * 60 + @second

      freeze # TimeOfDay instances are value objects
    end

    def <=>(other)
      return unless other.respond_to?(:second_of_day)
      @second_of_day <=> other.second_of_day
    end

    # Rounding to the given nearest number of seconds
    def round(round_sec = 1)
      down = self - (self.to_i % round_sec)
      up = down + round_sec

      difference_down = self - down
      difference_up = up - self

      if (difference_down < difference_up)
        return down
      else
        return up
      end
    end

    # Formats identically to Time#strftime
    def strftime(format_string)
      # Special case 2400 because strftime will load TimeOfDay into Time which
      # will convert 24 to 0
      format_string = format_string.gsub(/%H|%k/, '24') if @hour == 24
      Time.local(2000,1,1, @hour, @minute, @second).strftime(format_string)
    end

    def to_formatted_s(format = :default)
      if formatter = FORMATS[format]
        if formatter.respond_to?(:call)
          formatter.call(self).to_s
        else
          strftime(formatter)
        end
      else
        strftime "%H:%M:%S"
      end
    end
    alias_method :to_s, :to_formatted_s

    def value_for_database
      to_s
    end

    # Return a new TimeOfDay num_seconds greater than self. It will wrap around
    # at midnight.
    def +(num_seconds)
      TimeOfDay.from_second_of_day @second_of_day + num_seconds
    end

    # Return a new TimeOfDay num_seconds less than self. It will wrap around
    # at midnight.
    def -(other)
      if other.instance_of?(TimeOfDay)
        TimeOfDay.from_second_of_day @second_of_day - other.second_of_day
      else
        TimeOfDay.from_second_of_day @second_of_day - other
      end
    end

    # Returns a Time instance on date using self as the time of day
    # Optional time_zone will build time in that zone
    def on(date, time_zone=Tod::TimeOfDay.time_zone)
      time_zone.local date.year, date.month, date.day, @hour, @minute, @second
    end

    # Build a new TimeOfDay instance from second_of_day
    #
    #   TimeOfDay.from_second_of_day(3600) == TimeOfDay.new(1)   # => true
    def self.from_second_of_day(second_of_day)
      second_of_day = Integer(second_of_day)
      return new 24 if second_of_day == NUM_SECONDS_IN_DAY
      remaining_seconds = second_of_day % NUM_SECONDS_IN_DAY
      hour = remaining_seconds / NUM_SECONDS_IN_HOUR
      remaining_seconds -= hour * NUM_SECONDS_IN_HOUR
      minute = remaining_seconds / NUM_SECONDS_IN_MINUTE
      remaining_seconds -= minute * NUM_SECONDS_IN_MINUTE
      new hour, minute, remaining_seconds
    end
    class << self
      alias :from_i :from_second_of_day
    end

    # Build a TimeOfDay instance from string
    #
    # Strings only need to contain an hour. Minutes, seconds, AM or PM, and colons
    # are all optional.
    #   TimeOfDay.parse "8"                            # => 08:00:00
    #   TimeOfDay.parse "8am"                          # => 08:00:00
    #   TimeOfDay.parse "8pm"                          # => 20:00:00
    #   TimeOfDay.parse "8p"                           # => 20:00:00
    #   TimeOfDay.parse "9:30"                         # => 09:30:00
    #   TimeOfDay.parse "15:30"                        # => 15:30:00
    #   TimeOfDay.parse "3:30pm"                       # => 15:30:00
    #   TimeOfDay.parse "1230"                         # => 12:30:00
    #   TimeOfDay.parse "3:25:58"                      # => 03:25:58
    #   TimeOfDay.parse "515p"                         # => 17:15:00
    #   TimeOfDay.parse "151253"                       # => 15:12:53
    # You can give a block, that is called with the input if the string is not parsable.
    # If no block is given an ArgumentError is raised if try_parse returns nil.
    def self.parse(tod_string)
      try_parse(tod_string) || (block_given? ? yield(tod_string) : raise(ArgumentError, "Invalid time of day string"))
    end

    # Same as parse(), but return nil if not parsable (instead of raising an error)
    #   TimeOfDay.try_parse "8am"                      # => 08:00:00
    #   TimeOfDay.try_parse ""                         # => nil
    #   TimeOfDay.try_parse "abc"                      # => nil
    def self.try_parse(tod_string)
      tod_string = tod_string.to_s
      tod_string = tod_string.strip
      tod_string = tod_string.downcase
      tod_string = WORDS[tod_string] || tod_string
      if PARSE_24H_REGEX =~ tod_string || PARSE_12H_REGEX =~ tod_string
        hour, minute, second, a_or_p = $1.to_i, $2.to_i, $3.to_i, $4
        if hour == 12 && a_or_p == "a"
          hour = 0
        elsif hour < 12 && a_or_p == "p"
          hour += 12
        end

        new hour, minute, second
      else
        nil
      end
    end

    # Determine if a string is parsable into a TimeOfDay instance
    #   TimeOfDay.parsable? "8am"                      # => true
    #   TimeOfDay.parsable? "abc"                      # => false
    def self.parsable?(tod_string)
      !!try_parse(tod_string)
    end

    # If ActiveSupport TimeZone is available and set use current time zone else return Time
    def self.time_zone
      (Time.respond_to?(:zone) && Time.zone) || Time
    end

    def self.dump(time_of_day)
      time_of_day =
          if time_of_day.is_a? Hash
            # rails multiparam attribute
            # get hour, minute and second and construct new TimeOfDay object
            ::Tod::TimeOfDay.new(time_of_day[4], time_of_day[5], time_of_day[6])
          else
            # return nil, if input is not parsable
            Tod::TimeOfDay(time_of_day){}
          end
      time_of_day.to_s if time_of_day
    end

    def self.load(time)
      if time && !time.to_s.empty?
        return ::Tod::TimeOfDay.new(24) if time.respond_to?(:day) && time.day == 2 && time.hour == 0 && time.min == 0 && time.sec == 0
        ::Tod::TimeOfDay(time)
      end
    end
  end
end