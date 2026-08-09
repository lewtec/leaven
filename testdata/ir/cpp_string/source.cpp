// cpp_string — libstdc++ std::string (SSO append + size). stdio only.
// conda clang++ needs the bundled libstdc++ -isystem (see internal/genir).
#include <stdio.h>
#include <string>

int main() {
	std::string a = "foo";
	a += "bar";
	printf("%zu\n", a.size());
	return 0;
}
