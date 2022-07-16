
Detecting anycast addresses and more
===
Anycast networks are a pretty interesting way to fix quite a few issues with networked services that involve needing global spread.
One of the interesting things is that a computer cannot really tell (unless it has a full routing table of many providers in more than one geographic location maybe) that a IP is using anycast.

<h3>Anycast primer</h3>

Before I start this, I just want to go though what an anycast actually is.
Most programmers who have done networking have heard the comparison between a IP address and a telephone number, if you want to explain the idea of an anycast IP address to someone, I believe that the most simple way is to explain that if your computer is `34.24.77.2` (not my IP, not even routable) and lets say the direct translation to a phone number was `0121 441 526` and an anycast IP was `3.22.51.1`. Then its phone number would be like that of `999` (For Americans `911`, Other EU places `112`)
This kind of number is the "anycast" of telephone systems, because that phone number will be directed to the most *preferable* (hopefully the closest) call center in your region, this is why if I was to dial `999` in London, I would most likely be sent to a call center in London to take my emergency call, rather to one in Scotland.
Networks do the same thing basically.
Here is a small drawing of a network setup:

![anycast with routers](https://blog.benjojo.co.uk/asset/Kx5XPA0nxq)

Here you see a group of routers, all have the same IP address, Lets pretend that my router does know about all 3 of them, Because how most routing tables are setup, the packets will be sent to the fastest route possible (though, that does mean the fastest route it has, that may not in all cases mean the closest, "fastest" can in some cases also mean "cheapest") this is achieved by the "metric" on the route. Assuming the route metrics have been setup correctly. Packets will flow to the right place.

<h3>Detecting anycast addresses</h3>

Since there is nothing special about anycast addresses, other than that they are addresses that are "advertised" in the internet routing table in more than one place. There is no way to look at an address (other than what I mentioned above about seeing the global routing in many places) and know its anycasted.
One way for a human to check that an address is anycasted it to just trace route it from more than one location, Here is an example with 8.8.8.8 (an anycasted DNS server)
Server on West Coast USA:

```
ben@storm:~$ mtr -rwc 15 8.8.8.8
HOST: storm Loss% Snt Last Avg Best Wrst StDev
1. 162.244.92.1 0.0% 15 0.6 6.3 0.6 59.1 15.0
2. 10.1.1.5 0.0% 15 0.6 12.3 0.5 90.0 23.7
3. any2ix.coresite.com 0.0% 15 8.0 8.4 7.9 12.3 1.1
4. 209.85.250.99 0.0% 15 8.2 11.5 8.1 42.3 9.1
5. google-public-dns-a.google.com 0.0% 15 8.3 9.7 7.9 18.7 3.5
```

Server in Amsterdam:

```
ben@Spitfire:~$ mtr -rwc 15 8.8.8.8
HOST: Spitfire Loss% Snt Last Avg Best Wrst StDev
1. 95.46.198.1 0.0% 15 1.6 8.3 0.5 68.9 17.2
2. 80ge.cr0-br2-br3.smartdc.rtd.i3d.net 0.0% 15 1.3 4.4 0.5 13.8 4.5
3. 30ge.ar0-cr0.nikhef.ams.i3d.net 0.0% 15 1.5 17.9 1.5 56.5 17.8
4. core1.ams.net.google.com 0.0% 15 2.6 2.9 2.0 4.5 0.8
5. 209.85.241.237 0.0% 15 3.8 4.1 2.2 14.2 3.1
6. google-public-dns-a.google.com 0.0% 15 2.4 3.0 2.0 5.7 0.9
```

Notice on both of those servers, the round trip time is less than 10ms on each? Now unless routers have invented electron teleportation (the bare minimum time in a single direction to get from those two servers is [42ms](http://www.wolframalpha.com/input/?i=distance+between+LAX+and+AMS)) this is an anycast address
<h3>Discovery of users destination in regards to an anycast IP</h3>
This is normally a hard task for providers who use anycast, since you cannot guess where providers are going to route, since you cannot see their routing table or how their routers are configured in terms of what they have their metrics set to. They can only assume, For example “French IP addresses will hit the French announcement of the IP address”.
This is not always true, this is amplified more in regions of bad connectivity that have lesser connectivity than other regions, because they will end up buying off carriers that go direct to regions like Europe and North America where connectivity is better.
Where a user ends up if they try and contact an anycast IP depends on who their ISP is buying their greater connectivity off.
However, The provider can use their own anycast network to test where the user will end up, assuming that the target replies back when sent something (ICMP Echo, TCP SYN) you can listen back for the response. You don’t need to worry about where on the anycast network you send the probe, all you need to do is listen back for the response.
This works fine, as long as the target itself not an anycast address (and you can detect that as I have written above, by sending out from many nodes at once and seeing if the responses land back at different places)
<h3>PoC</h3>
I have a set of 3 servers (One on West Coast USA, East Coast USA, and one in Luxembourg, EU) that are setup to be any casted, and have assembled an anycast network out of the three servers.
Here is it detecting an anycast IP:

![anycast detection](https://blog.benjojo.co.uk/asset/6KqFDdwCqr)

and here is it testing against a set of unicast destinations:
University of Sydney:

![SYD Uni](https://blog.benjojo.co.uk/asset/5LbbWMakBj)

StackExchange (They are based in NY)

![StackOverflow](https://blog.benjojo.co.uk/asset/zWaW0E9IqA)

and finally the BBC:

![BBC](https://blog.benjojo.co.uk/asset/gECqxEup0S)

Want to try the tool out for yourself? You can find mine here: https://anycatch.benjojo.co.uk
Or if you want to build your own in the case that you have an anycast network infra of your own, You can find the code here: https://github.com/benjojo/AnyCatch
If you are using it, or find anything interesting with my own tool, Do let me know!
