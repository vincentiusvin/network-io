stuffs to write about:

1. exceeding the kernel's write buffer. your write() gets blocked until the client calls read()
   cat /proc/sys/net/ipv4/tcp_rmem
   cat /proc/sys/net/ipv4/tcp_wmem

   check the current buffer length by inspectig the Send-Q and Recv-Q value of
   ss -tulpna

   I want to find an exact threshold for this rather than just a "big number" but that doesn't seem possible since modern kernels auto tunes the value. Or maybe just my skill issue
