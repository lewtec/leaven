// cpp_point — struct, ctor, method, stack object. stdio only (no iostream).
#include <stdio.h>

struct Point {
	int x;
	int y;
	Point(int x, int y) : x(x), y(y) {}
	int manhattan() const {
		int ax = x < 0 ? -x : x;
		int ay = y < 0 ? -y : y;
		return ax + ay;
	}
};

int main(void) {
	Point p(3, -4);
	printf("%d\n", p.manhattan());
	return 0;
}
