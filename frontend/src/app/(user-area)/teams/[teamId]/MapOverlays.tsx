import Image from 'next/image';
import { useCallback, useMemo, useState } from 'react';
import { Marker } from 'react-map-gl/maplibre';

import RegionIcon from '@/assets/aws/Architecture-Group-Icons_02072025/Region_32.svg';
import ServerIcon from '@/assets/aws/Architecture-Group-Icons_02072025/Server-contents_32.svg';

import { MapLocation, MapRect } from './map';
import { MapMarkerPopup } from './MapMarkerPopup';
import { isEqualSelection, useSelection } from './selection';
import { useAwsRegions } from '@/hooks';
import { CombinedReport } from '@/report';

interface Props {
    combinedReport: CombinedReport;
    highlightPrincipalId?: string | null;
    zoom: number;
}

const meanMarkerLocation = (markers: MarkerData[]): MapLocation => {
    const meanMercatorX = markers.reduce((acc, marker) => acc + marker.location.mercatorX, 0) / markers.length;
    const meanMercatorY = markers.reduce((acc, marker) => acc + marker.location.mercatorY, 0) / markers.length;
    return MapLocation.fromMercator(meanMercatorX, meanMercatorY);
};

export type MarkerData =
    | {
          type: 'aws-region';
          location: MapLocation;
          id: string;
          name: string;
      }
    | {
          type: 'network';
          location: MapLocation;
          cidr: string;
      }
    | {
          type: 'cluster';
          location: MapLocation;
          rect: MapRect;
          awsRegionIds: Set<string>;
          networkCidrs: Set<string>;
          markerBounds: MapRect;
      };

const MARKER_ICON_SIZE = 22;

interface MarkerContentProps {
    data: MarkerData;
    fade?: boolean;
    emphasize?: boolean;
    statusText: string;
}

const MarkerContentInner = ({ data, fade, emphasize }: MarkerContentProps) => {
    const commonClasses = `h-[36px] min-w-[36px] ${emphasize ? 'border-2 border-white scale-120' : 'border-1 border-platinum'} ${fade ? 'opacity-20' : ''} rounded-full flex items-center justify-center transition-all duration-200 ease-in-out`;

    switch (data.type) {
        case 'aws-region':
            return (
                <div className={`bg-[#00A4A6] ${commonClasses}`}>
                    <Image src={RegionIcon} width={MARKER_ICON_SIZE} height={MARKER_ICON_SIZE} alt={data.name} />
                </div>
            );
        case 'network':
            return (
                <div className={`bg-[#7D8998] ${commonClasses}`}>
                    <Image src={ServerIcon} width={MARKER_ICON_SIZE} height={MARKER_ICON_SIZE} alt={data.cidr} />
                </div>
            );
        case 'cluster':
            let markerTypeClasses = 'bg-linear-to-br from-[#00A4A6] to-[#7D8998] from-50% to-50%';
            if (!data.networkCidrs.size) {
                markerTypeClasses = 'bg-[#00A4A6]';
            } else if (!data.awsRegionIds.size) {
                markerTypeClasses = 'bg-[#7D8998]';
            }
            return (
                <div className={`${markerTypeClasses} text-snow ${commonClasses}`}>
                    <span className="font-semibold text-xs">
                        {(data.awsRegionIds.size + data.networkCidrs.size).toLocaleString()}
                    </span>
                </div>
            );
        default:
            return null;
    }
};

const MarkerContent = (props: MarkerContentProps) => {
    const [hasMouse, setHasMouse] = useState(false);

    const onMouseEnter = useCallback(() => {
        setHasMouse(true);
    }, [setHasMouse]);

    const onMouseLeave = useCallback(() => {
        setHasMouse(false);
    }, [setHasMouse]);

    return (
        <div className="relative flex items-center justify-center">
            {hasMouse && <MapMarkerPopup data={props.data} statusText={props.statusText} />}
            <div onMouseEnter={onMouseEnter} onMouseLeave={onMouseLeave}>
                <MarkerContentInner {...props} />
            </div>
        </div>
    );
};

const clusterCentroidsForZoom = (zoom: number): number => 18 << Math.round(Math.log2(zoom + 1));

const minZoomForClusterCentroids = (centroids: number): number => {
    return Math.ceil(Math.pow(2, Math.log2(centroids / 18) - 0.5) - 1);
};

const maxZoomForClusterCentroids = (centroids: number): number => {
    return Math.floor(Math.pow(2, Math.log2(centroids / 18) + 0.5) - 1);
};

export const minZoomForClusterRect = (rect: MapRect): number => {
    const d = rect.maxMercatorX - rect.minMercatorX;
    return minZoomForClusterCentroids(Math.round(1.0 / d));
};

export const maxZoomForClusterRect = (rect: MapRect): number => {
    const d = rect.maxMercatorX - rect.minMercatorX;
    return maxZoomForClusterCentroids(Math.round(1.0 / d));
};

const clusterMarkers = (markers: MarkerData[], centroids: number): MarkerData[] => {
    const centroidRadius = 0.5 / centroids;
    const centroid = (x: number) => {
        return x - (x % (centroidRadius * 2)) + centroidRadius;
    };

    const clusters: {
        [key: string]: {
            awsRegionIds: Set<string>;
            networkCidrs: Set<string>;
            markers: MarkerData[];
            rect: MapRect;
            markerBounds: MapRect;
        };
    } = {};

    markers.forEach((marker) => {
        const x = centroid(marker.location.mercatorX);
        const y = centroid(marker.location.mercatorY);
        const key = `${x},${y}`;
        let cluster = clusters[key];
        if (!cluster) {
            cluster = {
                markers: [],
                awsRegionIds: new Set(),
                networkCidrs: new Set(),
                rect: new MapRect({
                    minMercatorX: x - centroidRadius,
                    minMercatorY: y - centroidRadius,
                    maxMercatorX: x + centroidRadius,
                    maxMercatorY: y + centroidRadius,
                }),
                markerBounds: new MapRect({
                    minMercatorX: marker.location.mercatorX,
                    minMercatorY: marker.location.mercatorY,
                    maxMercatorX: marker.location.mercatorX,
                    maxMercatorY: marker.location.mercatorY,
                }),
            };
            clusters[key] = cluster;
        }
        cluster.markers.push(marker);
        cluster.markerBounds.minMercatorX = Math.min(cluster.markerBounds.minMercatorX, marker.location.mercatorX);
        cluster.markerBounds.minMercatorY = Math.min(cluster.markerBounds.minMercatorY, marker.location.mercatorY);
        cluster.markerBounds.maxMercatorX = Math.max(cluster.markerBounds.maxMercatorX, marker.location.mercatorX);
        cluster.markerBounds.maxMercatorY = Math.max(cluster.markerBounds.maxMercatorY, marker.location.mercatorY);
        if (marker.type === 'aws-region') {
            cluster.awsRegionIds.add(marker.id);
        } else if (marker.type === 'network') {
            cluster.networkCidrs.add(marker.cidr);
        }
    });

    return Object.values(clusters).map((cluster) =>
        cluster.markers.length === 1
            ? cluster.markers[0]
            : {
                  type: 'cluster',
                  location: meanMarkerLocation(cluster.markers),
                  awsRegionIds: cluster.awsRegionIds,
                  networkCidrs: cluster.networkCidrs,
                  rect: cluster.rect,
                  markerBounds: cluster.markerBounds,
              },
    );
};

const markerHasNetwork = (marker: MarkerData, cidr: string): boolean => {
    if (marker.type === 'network') {
        return marker.cidr === cidr;
    }
    if (marker.type === 'cluster') {
        return marker.networkCidrs.has(cidr);
    }
    return false;
};

const markerHasAnyNetwork = (marker: MarkerData, networks: Set<string>): boolean => {
    if (marker.type === 'network') {
        return networks.has(marker.cidr);
    }
    if (marker.type === 'cluster') {
        return !marker.networkCidrs.isDisjointFrom(networks);
    }
    return false;
};

const markerHasAwsRegion = (marker: MarkerData, regionId: string): boolean => {
    if (marker.type === 'aws-region') {
        return marker.id === regionId;
    }
    if (marker.type === 'cluster') {
        return marker.awsRegionIds.has(regionId);
    }
    return false;
};

const markerHasAnyAwsRegion = (marker: MarkerData, regions: Set<string>): boolean => {
    if (marker.type === 'aws-region') {
        return regions.has(marker.id);
    }
    if (marker.type === 'cluster') {
        return !marker.awsRegionIds.isDisjointFrom(regions);
    }
    return false;
};

export const MapOverlays = (props: Props) => {
    const awsRegions = useAwsRegions();
    const [selection, setSelection] = useSelection();

    const { combinedReport, highlightPrincipalId, zoom } = props;

    const individualMarkers = useMemo<MarkerData[]>(() => {
        const markers: MarkerData[] = [];

        awsRegions.forEach((region) => {
            if (combinedReport.awsRegionIds.has(region.id)) {
                markers.push({
                    type: 'aws-region',
                    location: MapLocation.fromLatitudeAndLongitude(region.latitude, region.longitude),
                    id: region.id,
                    name: region.name,
                });
            }
        });

        Object.entries(combinedReport.networkLocations).forEach(([cidr, location]) => {
            markers.push({
                type: 'network',
                location: MapLocation.fromLatitudeAndLongitude(location.latitude, location.longitude),
                cidr,
            });
        });

        return markers;
    }, [awsRegions, combinedReport]);

    const clusterCentroids = clusterCentroidsForZoom(zoom);
    const markers = useMemo<MarkerData[]>(() => {
        return clusterMarkers(individualMarkers, clusterCentroids);
    }, [individualMarkers, clusterCentroids]);

    const hasSingleMarkerWithinSelection = useMemo<boolean>(() => {
        if (!selection || selection.type !== 'cluster') {
            return false;
        }
        let count = 0;
        for (const marker of markers) {
            if (
                marker.location.mercatorX >= selection.rect.minMercatorX &&
                marker.location.mercatorX <= selection.rect.maxMercatorX
            ) {
                if (
                    marker.location.mercatorY >= selection.rect.minMercatorY &&
                    marker.location.mercatorY <= selection.rect.maxMercatorY
                ) {
                    count++;
                    if (count > 1) {
                        return false;
                    }
                }
            }
        }
        return count === 1;
    }, [markers, selection]);

    const highlightNetworkCidrs = highlightPrincipalId && combinedReport.principals[highlightPrincipalId]?.networkCidrs;
    const highlightAwsRegions = highlightPrincipalId && combinedReport.principals[highlightPrincipalId]?.awsRegionIds;
    const selectedPrincipalNetworkCidrs =
        selection?.type === 'principal' && combinedReport.principals[selection.id]?.networkCidrs;
    const selectedPrincipalAwsRegions =
        selection?.type === 'principal' && combinedReport.principals[selection.id]?.awsRegionIds;

    return (
        <>
            {markers.map((marker) => {
                const isEquivalentClusterSelection =
                    hasSingleMarkerWithinSelection &&
                    selection?.type === 'cluster' &&
                    marker.type === 'cluster' &&
                    selection.rect.contains(marker.markerBounds);
                const isSelected = isEquivalentClusterSelection || isEqualSelection(marker, selection);

                const containsSelectedNetwork =
                    selection?.type === 'network' && markerHasNetwork(marker, selection.cidr);
                const containsSelectedAwsRegion =
                    selection?.type === 'aws-region' && markerHasAwsRegion(marker, selection.id);
                const containsSelectedCluster =
                    selection?.type === 'cluster' && marker.type === 'cluster' && marker.rect.contains(selection.rect);
                const containsSelectedPrincipalNetwork =
                    selection?.type === 'principal' &&
                    selectedPrincipalNetworkCidrs &&
                    markerHasAnyNetwork(marker, selectedPrincipalNetworkCidrs);
                const containsSelectedPrincipalAwsRegion =
                    selection?.type === 'principal' &&
                    selectedPrincipalAwsRegions &&
                    markerHasAnyAwsRegion(marker, selectedPrincipalAwsRegions);
                const containsSelectedPrincipalActivity =
                    containsSelectedPrincipalNetwork || containsSelectedPrincipalAwsRegion;
                const containsSelection =
                    containsSelectedNetwork ||
                    containsSelectedAwsRegion ||
                    containsSelectedCluster ||
                    containsSelectedPrincipalActivity;

                const containsHighlightedPrincipalNetwork =
                    highlightPrincipalId && highlightNetworkCidrs && markerHasAnyNetwork(marker, highlightNetworkCidrs);
                const containsHighlightedPrincipalAwsRegion =
                    highlightPrincipalId && highlightAwsRegions && markerHasAnyAwsRegion(marker, highlightAwsRegions);

                const shouldEmphasize = isSelected || containsSelection;
                const shouldFade =
                    !!highlightPrincipalId &&
                    !containsHighlightedPrincipalNetwork &&
                    !containsHighlightedPrincipalAwsRegion;

                const statusText = isSelected
                    ? 'Currently selected'
                    : containsSelectedPrincipalActivity
                      ? 'Contains selection activity'
                      : containsSelection
                        ? 'Contains selection'
                        : 'Click for more information';

                const extraClasses = shouldFade ? ' z-10' : shouldEmphasize ? ' z-100' : ' z-50';

                // XXX: we add the extra classes to the key to work around this bug:
                // https://github.com/visgl/react-map-gl/issues/2465
                // TODO: remove this when the bug fix is released
                const key = `${marker.location.mercatorX},${marker.location.mercatorY}` + extraClasses;

                return (
                    <Marker
                        key={key}
                        latitude={marker.location.latitude}
                        longitude={marker.location.longitude}
                        onClick={(e) => {
                            setSelection(isSelected ? null : marker);
                            e.originalEvent.preventDefault();
                            e.originalEvent.stopPropagation();
                        }}
                        className={
                            'hover:z-200 transition-all duration-200 ease-in-out baz cursor-pointer' + extraClasses
                        }
                    >
                        <MarkerContent
                            data={marker}
                            emphasize={shouldEmphasize}
                            fade={shouldFade}
                            statusText={statusText}
                        />
                    </Marker>
                );
            })}
        </>
    );
};
