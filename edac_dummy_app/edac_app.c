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

#define MEMORY_EDAC_FAIL -1

void sigalrm_handler( int sig )
{
	printf(" >>> Watchdog Expired \n");
	/* TODO: replace the printf with a call to
	 * the safe state function               */
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

	char *fileName = "/sys/devices/system/edac/mc/mc0/ue_count";

	if (pData != NULL) {
		fd = open(fileName, O_RDONLY);
		if (fd < 0)
			printf(" >>> Failed to open %s\n", fileName);
			/* TODO: replace the printf with a call to
			 * the safe state function               */
		else {
			status = read(fd, pData, lenBytes);
			close(fd);
		}
	}
	return status;
}


int main(int argc, char **argv)
{
    uint64_t ue_count;
    struct sigaction sact;
    int num_sent = 0;
    sigemptyset(&sact.sa_mask);
    sact.sa_flags = 0;
    sact.sa_handler = sigalrm_handler;
    sigaction(SIGALRM, &sact, NULL);

    ualarm(4000, 0);  /* Request SIGALRM in 4msec */
    while (1) {
    	MemoryGetUeCount(&ue_count);
    	if (ue_count) {
    		printf(">>>>> detected %lu uncorrectable errors\n", ue_count);
		/* TODO: replace the printf with a call to
		 * the safe state function               */
    	}
    	ualarm(4000, 0); /* Clear the previous counter and request SIGALRM in 4msec */
        usleep(1000); /* wait 1msec */
    }
    exit(0);
}
