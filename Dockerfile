FROM scratch
COPY octool /bin/octool
ENTRYPOINT [ "/bin/octool" ]
