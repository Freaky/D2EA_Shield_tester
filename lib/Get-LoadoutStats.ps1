Function Get-LoadoutStats{
    Param(
        $ShieldGenratorVariant,
        $ShieldBoosterLoadout,
        $ShieldBoosterVariantList
    )

    $ExpModifier = 1
    $KinModifier = 1
    $ThermModifier = 1
    $HitPointBonus = 0

    # Compute non diminishing returns modifiers
    ForEach($Booster in $ShieldBoosterLoadout){
        
        
        $BoosterStats = $($ShieldBoosterVariantList | Where-Object{$_.ID -eq $Booster})
        
        $ExpModifier = $ExpModifier * (1 - $BoosterStats.ExpResBonus)
        $KinModifier = $KinModifier * (1 - $BoosterStats.KinResBonus)
        $ThermModifier = $ThermModifier * (1 - $BoosterStats.ThermResBonus)

        $HitPointBonus = $HitPointBonus + $BoosterStats.ShieldStrengthBonus
    }

    # Compensate for diminishing returns
    If($ExpModifier -lt 0.7){
        $ExpModifier = 0.7 - (0.7-$ExpModifier)/2
    }
    If($KinModifier -lt 0.7){
        $KinModifier = 0.7 - (0.7-$KinModifier)/2
    }
    If($ThermModifier -lt 0.7){
        $ThermModifier = 0.7 - (0.7-$ThermModifier)/2
    }

    # Compute final Resistance
    $ExpRes = 1 - ((1 - $ShieldGenratorVariant.ExpRes) * $ExpModifier)
    $KinRes = 1 - ((1 - $ShieldGenratorVariant.KinRes) * $KinModifier)
    $ThermRes = 1 - ((1 - $ShieldGenratorVariant.ThermRes) * $ThermModifier)

    # Compute final Hitpoints

    $HitPoints = (1 + $HitPointBonus) * $ShieldGenratorVariant.ShieldStrength

    $LoadoutStat = New-Object PSCustomObject -Property @{
        HitPoints = [Double]$HitPoints + [Double]$SCBHitPoint + [Double]$GuardianShieldHitPoint
        RegenRate = [Double]$ShieldGenratorVariant.RegenRate
        ExplosiveResistance = [Double]$ExpRes
        KineticResistance = [Double]$KinRes
        ThermalResistance = [Double]$ThermRes
    }

    Return $LoadoutStat

   
}