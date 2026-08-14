// map<const V*, unsigned>: insert + iterate + virtual call on the key.
// csmith iv_bounds / Variable::loose_match hits this path.
#include <map>
#include <stdio.h>

struct V {
	virtual int id() const { return 42; }
	virtual bool same(const V *o) const { return this == o; }
};

static bool touches(const std::map<const V *, unsigned> &m, const V *v) {
	for (auto it = m.begin(); it != m.end(); ++it) {
		if (v->same(it->first)) {
			return true;
		}
	}
	return false;
}

int main() {
	V a, b;
	const V *pa = &a;
	const V *pb = &b;
	std::map<const V *, unsigned> m;
	// lvalue keys, same as csmith iv_bounds[iv] = bound
	m[pa] = 1;
	m[pb] = 2;
	unsigned n = 0;
	unsigned sum = 0;
	for (auto it = m.begin(); it != m.end(); ++it) {
		if (it->first == nullptr) {
			printf("nil key\n");
			return 1;
		}
		sum += it->second + (unsigned)it->first->id();
		n++;
	}
	printf("%u %u\n", n, sum);
	if (!touches(m, pa) || !touches(m, pb)) {
		printf("miss\n");
		return 1;
	}

	// CGContext copy-ctor clones iv_bounds. A bad pair copy leaves nil keys.
	std::map<const V *, unsigned> copied = m;
	unsigned cn = 0;
	unsigned csum = 0;
	for (auto it = copied.begin(); it != copied.end(); ++it) {
		if (it->first == nullptr) {
			printf("nil key copy\n");
			return 1;
		}
		csum += it->second + (unsigned)it->first->id();
		cn++;
	}
	printf("%u %u\n", cn, csum);

	// Bitwise copy of an empty map (llvm.memcpy / Go struct assign).
	// dest.left still names src.header; begin must equal dest.end.
	std::map<const V *, unsigned> empty_src;
	std::map<const V *, unsigned> empty_dst;
	__builtin_memcpy((void *)&empty_dst, (void *)&empty_src, sizeof(empty_src));
	unsigned en = 0;
	for (auto it = empty_dst.begin(); it != empty_dst.end(); ++it) {
		if (it->first == nullptr) {
			printf("nil key memcpy\n");
			return 1;
		}
		en++;
	}
	printf("%u\n", en);
	return 0;
}
