package game

import "math"

const experiencePerLevel = 10000

func enemyLevelForEncounter(floor int, template EnemyTemplate, elite bool) int {
	level := max(1, floor)
	if elite {
		level++
	}
	if template.BossTier > 0 {
		level += template.BossTier
	}
	return level
}

func enemyXPProgress(enemy *Enemy, player *Player, areaLevel int) int {
	if enemy == nil || player == nil {
		return 0
	}

	reward := enemyXPThreatScore(enemy)
	reward = reward * xpLevelDifferenceFactorBP(player.Level, enemy.Level) / 10000
	reward = reward * xpAreaFactorBP(player.Level, areaLevel) / 10000
	reward = reward * xpGearFactorBP(player, areaLevel) / 10000
	reward = reward * xpEncounterFactorBP(enemy) / 10000
	reward = reward * 100 / xpLevelDenominator(player.Level)

	if enemy.Elite {
		reward += 10
	}
	if enemy.Template.BossTier > 0 {
		reward += 26 * enemy.Template.BossTier
	}

	return clamp(reward, 4, 320)
}

func enemyXPThreatScore(enemy *Enemy) int {
	if enemy == nil {
		return 0
	}

	threat := 20
	threat += enemy.Template.XPReward * 3
	threat += enemy.Level * 5
	threat += enemy.Template.MaxHP / 3
	threat += enemy.AttackPower() * 2
	threat += enemy.DefensePower() * 4
	threat += enemy.Template.BurstDamage / 2
	if enemy.Template.AttackStatusChance > 0 || enemy.Template.PoisonChance > 0 {
		threat += 10
	}
	if enemy.Template.BurstStatusTurns > 0 {
		threat += 10 + enemy.Template.BurstStatusPotency*4
	}
	return threat
}

func xpLevelDenominator(level int) int {
	level = max(1, level)
	return 128 + level*24 + level*level*6
}

func xpLevelDifferenceFactorBP(playerLevel int, enemyLevel int) int {
	playerLevel = max(1, playerLevel)
	enemyLevel = max(1, enemyLevel)

	safeBand := 3 + playerLevel/16
	diff := abs(playerLevel - enemyLevel)
	if diff <= safeBand {
		if enemyLevel > playerLevel {
			return min(11200, 10000+(enemyLevel-playerLevel)*250)
		}
		return 10000
	}

	excess := diff - safeBand
	numerator := math.Pow(float64(playerLevel+6), 1.4)
	denominator := numerator + math.Pow(float64(excess+2), 2.15)*1.7
	factor := numerator / denominator
	return clamp(int(factor*10000+0.5), 1800, 9800)
}

func xpAreaFactorBP(playerLevel int, areaLevel int) int {
	diff := areaLevel - playerLevel
	switch {
	case diff >= 5:
		return 11000
	case diff >= 0:
		return 10000 + diff*160
	case diff >= -2:
		return 10000 + diff*220
	default:
		return clamp(10000+diff*320, 6500, 10000)
	}
}

func xpGearFactorBP(player *Player, areaLevel int) int {
	if player == nil {
		return 10000
	}

	gearScore := player.AttackPower()*2 + player.DefensePower()*3 + player.MaxHP()/4
	expected := 18 + areaLevel*4 + player.Level*2
	diff := gearScore - expected
	if diff > 0 {
		return clamp(10000-diff*120, 7600, 10000)
	}
	return clamp(10000+(-diff)*55, 10000, 10800)
}

func xpEncounterFactorBP(enemy *Enemy) int {
	if enemy == nil {
		return 10000
	}

	factor := 10000
	if enemy.Template.CanOpenDoors {
		factor += 250
	}
	if enemy.Template.BurstRange > 0 {
		factor += 350
	}
	if enemy.Template.AttackStatusChance > 0 || enemy.Template.PoisonChance > 0 {
		factor += 220
	}
	if enemy.Elite {
		factor += 900
	}
	if enemy.Template.BossTier > 0 {
		factor += 2200 + enemy.Template.BossTier*400
	}
	return factor
}

func hydrateProgressionState(game *Game) {
	if game == nil || game.Player == nil {
		return
	}

	legacyXPFormat := false
	if game.Floor != nil {
		for _, enemy := range game.Floor.Enemies {
			if enemy == nil {
				continue
			}
			if enemy.Level <= 0 {
				enemy.Level = enemyLevelForEncounter(game.Floor.Level, enemy.Template, enemy.Elite)
				legacyXPFormat = true
			}
		}
	}

	if legacyXPFormat {
		migrateLegacyXPProgress(game.Player)
	}
	game.Player.XP = clamp(game.Player.XP, 0, game.Player.NextLevelXP()-1)
}

func migrateLegacyXPProgress(player *Player) {
	if player == nil {
		return
	}

	start := legacyLevelStartXP(player.Level)
	end := legacyNextLevelXP(player.Level)
	progress := clamp(player.XP-start, 0, max(0, end-start))
	span := max(1, end-start)
	player.XP = progress * experiencePerLevel / span
}

func legacyLevelStartXP(level int) int {
	if level <= 1 {
		return 0
	}
	return legacyNextLevelXP(level - 1)
}

func legacyNextLevelXP(level int) int {
	level = max(1, level)
	return 16 + level*16 + (level-1)*(level-1)
}
