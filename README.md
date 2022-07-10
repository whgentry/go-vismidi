# MIDI controlled LED Go application

## About
The original intention of this project is create a go application that will will the midi input from a device file and control strips of LEDs. At first this will support an 88 key piano with a single string of LEDs. 

To aid in development the LED control will be abstracted away into an interface. A simulated led control package will be created before needing to control read hardware
