package aesCbc

/**
aes-cbc-256位加密解密实现
*/

/* Fixed Data */
var (
	InCo  = []byte{0xB, 0xD, 0x9, 0xE} /* Inverse Coefficients */
	fbsub = [256]byte{}
	rbsub = [256]byte{}
	ptab  = [256]byte{}
	ltab  = [256]byte{}

	ftable       = [256]uint32{}
	rtable       = [256]uint32{}
	rco          = [30]uint32{}
	isTablesInit = false
)

/* rotates x one bit to the left */
func ROTL(x byte) byte {
	return (x >> 7) | (x << 1)
}

/* Rotates 32-bit word left by 1, 2 or 3 byte  */
func ROTL8(x uint32) uint32 {
	return (x << 8) | (x >> 24)
}

func ROTL16(x uint32) uint32 {
	return (x << 16) | (x >> 16)
}

func ROTL24(x uint32) uint32 {
	return (x << 24) | (x >> 8)
}

/* Parameter-dependent data */
func pack(b []byte) uint32 { /* pack bytes into a 32-bit Word */
	return (uint32(b[3]) << 24) | (uint32(b[2]) << 16) | ((uint32(b[1]) << 8) | uint32(b[0]))
}

func unpack(a uint32, b []byte) { /* unpack bytes from a word */
	b[0] = byte(a)
	b[1] = byte(a >> 8)
	b[2] = byte(a >> 16)
	b[3] = byte(a >> 24)
}

func xtime(a byte) byte {
	var b byte
	if (a & 0x80) > 0 {
		b = 0x1B
	} else {
		b = 0
	}
	a <<= 1
	a ^= b
	return a
}

func bmul(x, y byte) byte { /* x.y= AntiLog(Log(x) + Log(y)) */
	if x > 0 && y > 0 {
		return ptab[(int(ltab[x])+int(ltab[y]))%255]
	}
	return 0
}

func subByte(a uint32) uint32 {
	b := make([]byte, 4)
	unpack(a, b)
	b[0] = fbsub[b[0]]
	b[1] = fbsub[b[1]]
	b[2] = fbsub[b[2]]
	b[3] = fbsub[b[3]]
	return pack(b)
}

func product(x, y uint32) byte { /* dot product of two 4-byte arrays */
	xb := make([]byte, 4)
	yb := make([]byte, 4)
	unpack(x, xb)
	unpack(y, yb)
	return bmul(xb[0], yb[0]) ^
		bmul(xb[1], yb[1]) ^
		bmul(xb[2], yb[2]) ^
		bmul(xb[3], yb[3])
}

func invMixCol(x uint32) uint32 { /* matrix Multiplication */
	var y, m uint32
	b := make([]byte, 4)
	m = pack(InCo)
	b[3] = product(m, x)
	m = ROTL24(m)
	b[2] = product(m, x)
	m = ROTL24(m)
	b[1] = product(m, x)
	m = ROTL24(m)
	b[0] = product(m, x)
	y = pack(b)
	return y
}

func byteSub(x byte) byte {
	y := ptab[255-ltab[x]] /* multiplicative inverse */
	x = y
	x = ROTL(x)
	y ^= x
	x = ROTL(x)
	y ^= x
	x = ROTL(x)
	y ^= x
	x = ROTL(x)
	y ^= x
	y ^= 0x63
	return y
}

func mcryptRijndaelGentables() { /* generate tables */
	var (
		i int
		y byte
		b = make([]byte, 4)
	)

	/* use 3 as primitive root to generate power and log tables */
	ltab[0] = 0
	ptab[0] = 1
	ltab[1] = 0
	ptab[1] = 3
	ltab[3] = 1
	for i = 2; i < 256; i++ {
		ptab[i] = ptab[i-1] ^ xtime(ptab[i-1])
		ltab[ptab[i]] = byte(i)
	}

	/* affine transformation:- each bit is xored with itself shifted one bit */
	fbsub[0] = 0x63
	rbsub[0x63] = 0
	for i = 1; i < 256; i++ {
		y = byteSub(byte(i))
		fbsub[i] = y
		rbsub[y] = byte(i)
	}
	i = 0
	y = 1
	for ; i < 30; i++ {
		rco[i] = uint32(y)
		y = xtime(y)
	}

	/* calculate forward and reverse tables */
	for i = 0; i < 256; i++ {
		y = fbsub[i]
		b[3] = y ^ xtime(y)
		b[2] = y
		b[1] = y
		b[0] = xtime(y)
		ftable[i] = pack(b)

		y = rbsub[i]
		b[3] = bmul(InCo[0], y)
		b[2] = bmul(InCo[1], y)
		b[1] = bmul(InCo[2], y)
		b[0] = bmul(InCo[3], y)
		rtable[i] = pack(b)
	}
}

func mcryptSetKey(rinst *RI, key []byte, nk int) int { /* blocksize=32*nb bits. Key=32*nk bits */
	/* currently nb,bk = 4, 6 or 8          */
	/* key comes as 4*rinst.Nk bytes              */
	/* Key Scheduler. Create expanded encryption key */
	nb := 8 /* 256 block size */
	var (
		i, j, k, m, N int
		C1, C2, C3    int
		CipherKey     = make([]uint32, 8)
	)

	nk /= 4
	if !isTablesInit {
		mcryptRijndaelGentables()
		isTablesInit = true
	}

	rinst.Nb = nb
	rinst.Nk = nk

	/* rinst.Nr is number of rounds */
	if rinst.Nb >= rinst.Nk {
		rinst.Nr = 6 + rinst.Nb
	} else {
		rinst.Nr = 6 + rinst.Nk
	}

	C1 = 1
	if rinst.Nb < 8 {
		C2 = 2
		C3 = 3
	} else {
		C2 = 3
		C3 = 4
	}

	/* pre-calculate forward and reverse increments */
	m, j = 0, 0
	for ; j < nb; j++ {
		rinst.fi[m] = byte((j + C1) % nb)
		rinst.fi[m+1] = byte((j + C2) % nb)
		rinst.fi[m+2] = byte((j + C3) % nb)
		rinst.ri[m] = byte((nb + j - C1) % nb)
		rinst.ri[m+1] = byte((nb + j - C2) % nb)
		rinst.ri[m+2] = byte((nb + j - C3) % nb)
		m += 3
	}

	N = rinst.Nb * (rinst.Nr + 1)
	i, j = 0, 0
	for ; i < rinst.Nk; i++ {
		CipherKey[i] = pack(key[j:])
		j += 4
	}
	for i = 0; i < rinst.Nk; i++ {
		rinst.fkey[i] = CipherKey[i]
	}
	j = rinst.Nk
	k = 0
	for ; j < N; j += rinst.Nk {
		rinst.fkey[j] = rinst.fkey[j-rinst.Nk] ^
			subByte(ROTL24(rinst.fkey[j-1])) ^ rco[k]
		if rinst.Nk <= 6 {
			for i = 1; i < rinst.Nk && (i+j) < N; i++ {
				rinst.fkey[i+j] = rinst.fkey[i+j-rinst.Nk] ^ rinst.fkey[i+j-1]
			}
		} else {
			for i = 1; i < 4 && (i+j) < N; i++ {
				rinst.fkey[i+j] = rinst.fkey[i+j-rinst.Nk] ^ rinst.fkey[i+j-1]
			}
			if j+4 < N {
				rinst.fkey[j+4] = rinst.fkey[j+4-rinst.Nk] ^ subByte(rinst.fkey[j+3])
			}
			for i = 5; i < rinst.Nk && (i+j) < N; i++ {
				rinst.fkey[i+j] = rinst.fkey[i+j-rinst.Nk] ^ rinst.fkey[i+j-1]
			}
		}
		k++
	}

	/* now for the expanded decrypt key in reverse order */
	for j = 0; j < rinst.Nb; j++ {
		rinst.rkey[j+N-rinst.Nb] = rinst.fkey[j]
	}
	for i = rinst.Nb; i < N-rinst.Nb; i += rinst.Nb {
		k = N - rinst.Nb - i
		for j = 0; j < rinst.Nb; j++ {
			rinst.rkey[k+j] = invMixCol(rinst.fkey[i+j])
		}
	}
	for j = N - rinst.Nb; j < N; j++ {
		rinst.rkey[j-N+rinst.Nb] = rinst.fkey[j]
	}

	return 0
}

/* There is an obvious time/space trade-off possible here.     *
 * Instead of just one ftable[], I could have 4, the other     *
 * 3 pre-rotated to save the ROTL8, ROTL16 and ROTL24 overhead */
func mcryptEncrypt(rinst *RI, buff []byte) {
	var (
		i, j, k, m int
		a          = make([]uint32, 8)
		b          = make([]uint32, 8)
		x, y, t    []uint32
	)
	i, j = 0, 0
	for ; i < rinst.Nb; i++ {
		a[i] = pack(buff[j:])
		a[i] ^= rinst.fkey[i]
		j += 4
	}
	k = rinst.Nb
	x = a
	y = b

	/* State alternates between a and b */
	for i = 1; i < rinst.Nr; i++ { /* rinst.Nr is number of rounds. May be odd. */
		/* if rinst.Nb is fixed - unroll this next
		   loop and hard-code in the values of fi[]  */
		m, j = 0, 0
		for ; j < rinst.Nb; j++ { /* deal with each 32-bit element of the State */
			/* This is the time-critical bit */
			y[j] = rinst.fkey[k] ^ ftable[byte(x[j])] ^
				ROTL8(ftable[(byte)(x[rinst.fi[m]]>>8)]) ^
				ROTL16(ftable[(byte)(x[rinst.fi[m+1]]>>16)]) ^
				ROTL24(ftable[x[rinst.fi[m+2]]>>24])
			m += 3
			k++
		}
		t = x
		x = y
		y = t /* swap pointers */
	}

	/* Last Round - unroll if possible */
	m, j = 0, 0
	for ; j < rinst.Nb; j++ {
		y[j] = rinst.fkey[k] ^ uint32(fbsub[byte(x[j])]) ^
			ROTL8(uint32(fbsub[(byte)(x[rinst.fi[m]]>>8)])) ^
			ROTL16(uint32(fbsub[(byte)(x[rinst.fi[m+1]]>>16)])) ^
			ROTL24(uint32(fbsub[x[rinst.fi[m+2]]>>24]))
		m += 3
		k++
	}
	i, j = 0, 0
	for ; i < rinst.Nb; i++ {
		unpack(y[i], buff[j:])
		x[i], y[i] = 0, 0 /* clean up stack */
		j += 4
	}
}

func mcryptDecrypt(rinst *RI, buff []byte) {
	var (
		i, j, k, m int
		a          = make([]uint32, 8)
		b          = make([]uint32, 8)
		x, y, t    []uint32
	)
	i, j = 0, 0
	for ; i < rinst.Nb; i++ {
		a[i] = pack(buff[j:])
		a[i] ^= rinst.rkey[i]
		j += 4
	}
	k = rinst.Nb
	x = a
	y = b

	/* State alternates between a and b */
	for i = 1; i < rinst.Nr; i++ { /* rinst.Nr is number of rounds. May be odd. */
		/* if rinst.Nb is fixed - unroll this next
		   loop and hard-code in the values of ri[]  */
		m, j = 0, 0
		for ; j < rinst.Nb; j++ { /* This is the time-critical bit */
			y[j] = rinst.rkey[k] ^ rtable[byte(x[j])] ^
				ROTL8(rtable[(byte)(x[rinst.ri[m]]>>8)]) ^
				ROTL16(rtable[(byte)(x[rinst.ri[m+1]]>>16)]) ^
				ROTL24(rtable[x[rinst.ri[m+2]]>>24])
			m += 3
			k++
		}
		t = x
		x = y
		y = t /* swap pointers */
	}

	/* Last Round - unroll if possible */
	m, j = 0, 0
	for ; j < rinst.Nb; j++ {
		y[j] = rinst.rkey[k] ^ uint32(rbsub[byte(x[j])]) ^
			ROTL8(uint32(rbsub[(byte)(x[rinst.ri[m]]>>8)])) ^
			ROTL16(uint32(rbsub[(byte)(x[rinst.ri[m+1]]>>16)])) ^
			ROTL24(uint32(rbsub[x[rinst.ri[m+2]]>>24]))
		m += 3
		k++
	}
	i, j = 0, 0
	for ; i < rinst.Nb; i++ {
		unpack(y[i], buff[j:])
		j += 4
	}
}
