# How to Test

## Software Setup
Before you can test Speedgoat IO, you'll need MATLAB installed with Simulink and Simulink Real-time (preferably MATLAB 
R2023a). This is so you can run the corresponding server on the Speedgoat to allow it to read/send data back. For 
information on setting up this environment, see the Speedgoat quick start guide.

## Hardware Setup
You'll need to connect an Ethernet cable between your test machine and the Speedgoat's host link port. There isn't a
required OS to run the test, the test machine just needs to be able to connect via Ethernet.

## Running the Test
Start by connecting to the Speedgoat in Simulink Real-time and run the server (.slx file for this should be available on
GitHub eventually). It may take up to a few minutes to build and run on the Speedgoat.

Once running, execute the test file and monitor the displays in Simulink or the physical IO pins for the appropriate 
response. Digital pin 9 and 16 should be set high and analog output 1 should be set to 2.5V initially (refer to the I/O 
Connector sheet if needed). After 2 seconds, digital pin 16 should be set low and analog output 1 should be set to 3V.

If this does not occur, it may mean that the sample time for the TCP Receive block (in Simulink) is too high. Try 
setting it to 0.01 or lower.