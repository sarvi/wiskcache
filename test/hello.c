#include <stdio.h>
#include "sum.h"
int main()
{
   printf("Hello World\n");
   int a = 10, b = 10;
   int sum = add(a, b);
   printf("Addition of %d and %d = %d\n", a, b, sum);
   return 0;
}
