

In a busy IT world remembering to disable monitoring during ad-hoc system shutdowns or reboots can be challenging.  I thought I’d have a go at making an API call to the vROps (or Aria) from the affected virtual machines so they could call home in the event of a shutdown - set themselves as ‘in maintenance’ and then when they start back up un-do the maintenance and life would continue on as normal..

…no alarms, no grumpy staff (hopefully!) ;) and we all carry on as if nothing had ever happened.

I’ve made this as OS agnostic as possible so that you can compile this on the systems you require.

Please note that I’ve simply added the API login credentials to the code below, please use your preferred method to secure things according to your own internal practices.

This code is a starter for ten - feel free to modify this as you see fit.

Here’s a high level description of the codes workflow :

1. Prep the data required.
2. Grab an API auth token.
3. Grab the vRO objectID from the VM Name supplied in command line.
4. Set the objectID as maintained or delete the maintained flag depending on the command line used.

The program needs 3 arguments -  

1. vROps server fqdn.
2. VM name
3. start/end

You’ll need to get each operating system to execute the code (compiled or built as you see fit) during system startup and shutdown so that the corresponding vROps object is also silenced and un-silenced.

There are many many ways to achieve this for the various operating systems so i’ll leave that part up to you.

The Go code I’ve found is the better way to create a small binary that will execute on the windows systems required.




Now you can feel free to reboot all the things! (with permission of course!) :)
