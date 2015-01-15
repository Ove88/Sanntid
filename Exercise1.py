from threading import Thread
global i
i = 0

def someThreadFunction1():
	global i
	for s in range(0,1000000):
		i = i+1
def someThreadFunction2():
	global i
	for s in range(0,1000000):
		i = i-1


def main():
	someThread1 = Thread(target = someThreadFunction1, args = (),)
	someThread2 = Thread(target = someThreadFunction2, args = (),)
	someThread1.start()
	someThread2.start()
    
	someThread1.join()
	someThread2.join()
	print(i)

main()
