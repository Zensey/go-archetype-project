# fiatconv

## Usage: 

`fiatconv <amount_src:float> <src_symbol:string> <dst_symbol:string>`

Arguments:
*  amount_src  Amount to convert
*  src_symbol  Currency you are converting from
*  dst_symbol  Currency you are converting to

Example: 

`fiatconv 100 EUR GBP`

Output:

`100.00 EUR = 90.43 GBP`

## Makefile rules
* `make fiatconv`
* `make docker-build`
* `make lint`
