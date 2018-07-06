


def solu(A):
    bin=0
    a=[]
    for each in A:
        bin+=pow(2,each)
    print(bin)
    i=1
    while i<bin:
        print('A')
        i
    #     if(i**2>bin):
    #         a.append(i-1)
    #         bin=bin-i-1
    #         i=0
    # return len(a)

solu([1,0,2,0,0,2])
