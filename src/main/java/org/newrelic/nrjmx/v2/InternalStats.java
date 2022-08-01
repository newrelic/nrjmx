/*
 * Copyright 2021 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package org.newrelic.nrjmx.v2;

import org.newrelic.nrjmx.v2.nrprotocol.InternalStat;

import java.util.ArrayList;
import java.util.List;

/**
 * InternalStats class used to collect internal nrjmx query stats for troubleshooting.
 */
public class InternalStats {
    private List<InternalStat> stats;

    /* maxSize defines how many stats we can keep in memory. When limit is reached old ones are discarded. */
    private long maxSize;

    final static private long DEFAULT_MAX_SIZE = 100000;

    /**
     * InternalStats constructor.
     *
     * @param maxSize long defines how many stats we can keep in memory. When limit is reached old ones are discarded.
     */
    public InternalStats(long maxSize) {
        if (maxSize < 1) {
            maxSize = DEFAULT_MAX_SIZE;
        }

        this.maxSize = maxSize;
        this.stats = new ArrayList<>();
    }

    /**
     * Records a new InternalStat and returns it for attaching more data if required.
     *
     * @param statType String name of the stat.
     * @return InternalStat new registered stat to add data to it.
     */
    public InternalStat record(String statType) {
        synchronized (stats) {
            if (stats.size() >= maxSize) {
                stats.remove(0);
            }

            InternalStat stat = new InternalStat()
                    .setStartTimestamp(System.currentTimeMillis())
                    .setMilliseconds((double) System.nanoTime() / 1000000.0)
                    .setStatType(statType);
            stats.add(stat);
            return stat;
        }
    }

    /**
     * Returns all the collected stats and clear them.
     *
     *
     * @return List<InternalStat> returns all collected internal stats.
     */
    public List<InternalStat> getStats() {
        synchronized (this.stats) {
            List<InternalStat> stats = new ArrayList<>(this.stats);
            this.stats = new ArrayList<>();
            return stats;
        }
    }

    /**
     * Calculate the elapsed ms with .3f precision since the stat was recorded and attach it to the stat.
     */
    public static void setElapsedMs(InternalStat internalStat) {
        internalStat.setMilliseconds((double) System.nanoTime() / 1000000.0 - internalStat.milliseconds);
    }
}
