import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class SettingsService {
  private _visible = false;
  set visible(v: boolean) {
    this._visible = v;
    this.visibleSubject.next(v);
  }
  get visible(): boolean {
    return this._visible;
  }

  private visibleSubject = new BehaviorSubject<boolean>(false);

  get visible$(): Observable<boolean> {
    return this.visibleSubject.asObservable();
  }

  constructor() {}
}
