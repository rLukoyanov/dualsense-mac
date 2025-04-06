byte 0, bit 0: Report ID - 8 bits
    always 0x01
byte 1, bit 0: Generic Desktop / X - 8 bits - left stick X axis
    left: 0x00, right: 0xff, neutral: ~0x80
byte 2, bit 0: Generic Desktop / Y - 8 bits - left stick Y axis
    up: 0x00, down: 0xff, neutral: ~0x80
byte 3, bit 0: Generic Desktop / Z - 8 bits - right stick X axis
    left: 0x00, right: 0xff, neutral: ~0x80
byte 4, bit 0: Generic Desktop / Rz - 8 bits - right stick Y axis
    up: 0x00, down: 0xff, neutral: ~0x80
byte 5, bit 0: Generic Desktop / Hat switch - 4 bits - directional buttons
    neutral: 0x8, N: 0x0, NE: 0x1, E: 0x2, SE: 0x3, S: 0x4, SW: 0x5, W: 0x6, NW: 0x7
byte 5, bit 4: Button / 0x01 - 1 bit - Square button
byte 5, bit 5: Button / 0x02 - 1 bit - Cross button
byte 5, bit 6: Button / 0x03 - 1 bit - Circle button
byte 5, bit 7: Button / 0x04 - 1 bit - Triangle button
byte 6, bit 0: Button / 0x05 - 1 bit - L1 button
byte 6, bit 1: Button / 0x06 - 1 bit - R1 button
byte 6, bit 2: Button / 0x07 - 1 bit - L2 button
byte 6, bit 3: Button / 0x08 - 1 bit - R2 button
byte 6, bit 4: Button / 0x09 - 1 bit - Create button
byte 6, bit 5: Button / 0x0a - 1 bit - Options button
byte 6, bit 6: Button / 0x0b - 1 bit - L3 button
byte 6, bit 7: Button / 0x0c - 1 bit - R3 button
byte 7, bit 0: Button / 0x0d - 1 bit - PS button
byte 7, bit 1: Button / 0x0e - 1 bit - Touchpad button
byte 7, bit 2: Vendor defined 0xFF00 / 0x21 - 6 bits
byte 8, bit 0: Generic Desktop / Rx - 8 bits - L2 axis
    neutral: 0x00, pressed: 0xff
byte 9, bit 0: Generic Desktop / Ry - 8 bits - R2 axis
    neutral: 0x00, pressed: 0xff