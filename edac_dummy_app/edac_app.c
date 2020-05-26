// SPDX-License-Identifier: GPL-2.0
/*
 ============================================================================
 Name        : Edac_app.c
 Author      : gabriele.paoloni@intel.com
 Version     :
 Copyright   : Copyright (c) 2020, Intel Corporation.
 Description : Hello World in C, Ansi-style
 ============================================================================
 */

#include <signal.h>
#include <unistd.h>
#include <stdio.h>
#include <stdint.h>
#include <fcntl.h>
#include <unistd.h>
#include <stdlib.h>
#include <gpiod.h>

#define MEMORY_EDAC_FAIL -1
#define GPIO_LINE 1

struct gpiod_line *line;

void sigalrm_handler( int sig)
{
	printf(" >>> Watchdog Expired \n");
	/* Drive Safe State */
	gpiod_line_set_value(line, 1);
	exit(0);
}

/**
Get UE count from platform.
It gets the Un-correcatble error count for IBECC from the platform.

@param[out] *pData - pointer to data read

@retval the number of bytes read - if succeed
@retval -1 - if failed for whatever reason (and errno is set appropriately)
 **/
uint32_t MemoryGetUeCount(uint64_t *pData)
{
	uint32_t status = MEMORY_EDAC_FAIL;
	int32_t fd;
	uint32_t lenBytes = 4;
	char buf[4];

	char *fileName = "/sys/devices/system/edac/mc/mc0/ue_count";

	if (pData != NULL) {
		fd = open(fileName, O_RDONLY);
		if (fd < 0) {
			printf(" >>> Failed to open %s\n", fileName);
			/* Drive Safe State */
			gpiod_line_set_value(line, 1);
		} else {
			status = read(fd, buf, lenBytes);
			*pData = atol(buf);
			printf(" >>> pData is %ld\n", *pData);
			close(fd);
		}
	}
	return status;
}


int main(int argc, char **argv)
{
	uint64_t ue_count = 0;
	struct sigaction sact;
	int num_sent = 0;
	int ret = -1;
	struct gpiod_chip *chip;
	struct gpiod_line_bulk bulk;
	int offset = 0;
	int gpio_val = 0;
	int gpio_read = 0;

	char gpio_chip_name[] = "/dev/gpiochip0";

	chip = gpiod_chip_open(gpio_chip_name);
	if (!chip) {
		printf(">>> gpiod_chip_open failed\n");
		return -1;
	}

	/*TODO: for some reason gpiod_chip_get_all_lines() segfaults so
	 * it needs to be investigated. The code below is a workaround */
	do {
		line = gpiod_chip_get_line(chip, offset++);
		if ((line) && !gpiod_line_is_used(line)) {
			ret = gpiod_line_request_output(line, "edac-diagnostic",
								GPIOD_LINE_ACTIVE_STATE_LOW);
			if (!ret)
				break;
		}
	} while (line);

	if (ret) {
		printf("Unable to find a free IRQ line\n");
		gpiod_chip_close(chip);
		return -1;
	}
	printf("using GPIO line %d", (offset-1));

	sigemptyset(&sact.sa_mask);
	sact.sa_flags = 0;
	sact.sa_handler = sigalrm_handler;
	sigaction(SIGALRM, &sact, NULL);

	ualarm(4000, 0);  /* Request SIGALRM in 4msec */
	while (1) {
		gpiod_line_set_value(line, 0);
		ret = MemoryGetUeCount(&ue_count);
		if (ret < 0)
			break;
		if (ue_count) {
			printf(">>>>> detected %lu uncorrectable errors\n", ue_count);
			/* Drive Safe State */
			gpiod_line_set_value(line, 1);
			break;
		}
		ualarm(4000, 0); /* Clear the previous counter and request SIGALRM in 4msec */
		usleep(1000); /* wait 1msec */
	}
	exit(0);
}
