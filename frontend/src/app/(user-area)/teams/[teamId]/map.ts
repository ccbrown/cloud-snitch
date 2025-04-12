const mercatorXfromLng = (lng: number) => {
    return (180 + lng) / 360;
};

const mercatorYfromLat = (lat: number) => {
    return (180 - (180 / Math.PI) * Math.log(Math.tan(Math.PI / 4 + (lat * Math.PI) / 360))) / 360;
};

const lngFromMercatorX = (x: number) => {
    return x * 360 - 180;
};

const latFromMercatorY = (y: number) => {
    const y2 = 180 - y * 360;
    return (360 / Math.PI) * Math.atan(Math.exp((y2 * Math.PI) / 180)) - 90;
};

export class MapLocation {
    latitude: number;
    longitude: number;
    mercatorX: number;
    mercatorY: number;

    constructor(params: { latitude: number; longitude: number; mercatorX: number; mercatorY: number }) {
        this.latitude = params.latitude;
        this.longitude = params.longitude;
        this.mercatorX = params.mercatorX;
        this.mercatorY = params.mercatorY;
    }

    static fromLatitudeAndLongitude(latitude: number, longitude: number): MapLocation {
        const mercatorX = mercatorXfromLng(longitude);
        const mercatorY = mercatorYfromLat(latitude);
        return new MapLocation({
            latitude,
            longitude,
            mercatorX,
            mercatorY,
        });
    }

    static fromMercator(mercatorX: number, mercatorY: number): MapLocation {
        return new MapLocation({
            latitude: latFromMercatorY(mercatorY),
            longitude: lngFromMercatorX(mercatorX),
            mercatorX,
            mercatorY,
        });
    }

    toString(): string {
        const latDirection = this.latitude >= 0 ? 'N' : 'S';
        const lat = Math.abs(this.latitude).toFixed(4);
        const lonDirection = this.longitude >= 0 ? 'E' : 'W';
        const lon = Math.abs(this.longitude).toFixed(4);
        return `${lat}° ${latDirection}, ${lon}° ${lonDirection}`;
    }
}

export class MapRect {
    minMercatorX: number;
    minMercatorY: number;
    maxMercatorX: number;
    maxMercatorY: number;

    constructor(params: { minMercatorX: number; minMercatorY: number; maxMercatorX: number; maxMercatorY: number }) {
        this.minMercatorX = params.minMercatorX;
        this.minMercatorY = params.minMercatorY;
        this.maxMercatorX = params.maxMercatorX;
        this.maxMercatorY = params.maxMercatorY;
    }

    contains(other: MapRect): boolean {
        return (
            this.minMercatorX <= other.minMercatorX &&
            this.minMercatorY <= other.minMercatorY &&
            this.maxMercatorX >= other.maxMercatorX &&
            this.maxMercatorY >= other.maxMercatorY
        );
    }

    containsLocation(location: MapLocation): boolean {
        return (
            this.minMercatorX <= location.mercatorX &&
            this.minMercatorY <= location.mercatorY &&
            this.maxMercatorX >= location.mercatorX &&
            this.maxMercatorY >= location.mercatorY
        );
    }

    equals(other: MapRect): boolean {
        return (
            this.minMercatorX === other.minMercatorX &&
            this.minMercatorY === other.minMercatorY &&
            this.maxMercatorX === other.maxMercatorX &&
            this.maxMercatorY === other.maxMercatorY
        );
    }
}
