# Game

## run server
`./game`
 
  
## run client: 
`./game -game_server 85.253.108.246:34259`

## Notes
Under some circumstances, lets say you have

* `NAT mapping behavior: endpoint independent`
* `NAT filtering behavior: address and port dependent`

to make hole punching work will take you to provide
a server with the client's public adress. For that purpose
input the client's public address onto the server console.

## Further improvements

To improve the usability like in the case described above
a peer broker can be introduced. In such architecture the flow will be different:
* server registers its address in broker and gets an uid.
* uid is announced
* client registers its <address, uid> in broker and gets server address in return 
* broker propagates peer address to server
 