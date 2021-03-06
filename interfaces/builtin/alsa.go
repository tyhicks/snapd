// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2016 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package builtin

const alsaDescription = `
The alsa interface allows connected plugs to access raw ALSA devices.

The core snap provides the slot that is shared by all the snaps.
`

const alsaConnectedPlugAppArmor = `
# Description: Allow access to raw ALSA devices.

/dev/snd/  r,
/dev/snd/* rw,

/run/udev/data/c116:[0-9]* r, # alsa

# Allow access to the alsa state dir
/var/lib/alsa/{,*}         r,
`

func init() {
	registerIface(&commonInterface{
		name:                  "alsa",
		description:           alsaDescription,
		connectedPlugAppArmor: alsaConnectedPlugAppArmor,
		reservedForOS:         true,
	})
}
